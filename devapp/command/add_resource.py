# -*- coding: utf-8 -*-

import os
import json
import shutil
import sys
import time
import requests
import zipfile

from command_base import BaseCommand

code_zip_path = './codebase.zip'
TEMPLATE_FILE_DIR = '_gofile_template'

class Command(BaseCommand):
	help = "add_resource [package] [resource]"
	args = ''

	def confirm_dir_exists(self, dir):
		if not os.path.exists(dir):
			print 'make dir: ', dir
			os.makedirs(dir)
		else:
			print 'dir(%s) is already exists' % dir

	def check_file_exists(self, file_path):
		return os.path.exists(file_path)

	def get_plural_name(self, name):
		plural_name = name
		if plural_name[-1] == 'y':
			plural_name = '%sies' % plural_name[:-1]
		elif plural_name[-1] == 's':
			plural_name = '%ses' % plural_name
		else:
			plural_name = '%ss' % plural_name

		return plural_name

	def get_names(self, resource):
		#make class name
		items = resource.split('_')
		for i, item in enumerate(items):
			items[i] = item.capitalize()
		class_name = ''.join(items)

		#make plural class name
		plural_class_name = self.get_plural_name(class_name)

		#make plural resource name
		plural_resource = self.get_plural_name(resource)

		#get service name
		with open('./service.json', 'rb') as f:
			content = f.read()
			service_info = json.loads(content)

		return {
			'service_name': service_info['name'],
			'model_name': class_name,
			'resource_class_name': class_name,
			'entity_class_name': class_name,
			'plural_class_name': plural_class_name,
			'plural_resource': plural_resource
		}

	def generate_model_file(self, names):
		print '\n>>>>>>>>>> generate model <<<<<<<<<<'
		dir_path = '_generate/models'
		self.confirm_dir_exists(dir_path)
		file_path = os.path.join(dir_path, '%(resource)s.go' % names)
		
		self.render_file_to('model.go', file_path, names)

	def generate_resource_file(self, names):
		print '\n>>>>>>>>>> generate rest resource <<<<<<<<<<'
		dir_path = '_generate/rest'
		self.confirm_dir_exists(dir_path)

		file_path = os.path.join(dir_path, '%(resource)s.go' % names)
		self.render_file_to('resource.go', file_path, names)
		
		file_path = os.path.join(dir_path, '%(plural_resource)s.go' % names)
		self.render_file_to('resources.go', file_path, names)

	def generate_business_file(self, names):
		print '\n>>>>>>>>>> generate business objects <<<<<<<<<<'
		dir_path = '_generate/business'
		self.confirm_dir_exists(dir_path)

		file_path = os.path.join(dir_path, '%(resource)s.go' % names)
		self.render_file_to('entity.go', file_path, names)
		
		file_path = os.path.join(dir_path, '%(resource)s_repository.go' % names)
		self.render_file_to('repository.go', file_path, names)

		file_path = os.path.join(dir_path, 'resp_%(package)s.go' % names)
		self.render_file_to('resp.go', file_path, names)

		file_path = os.path.join(dir_path, 'encode_%(resource)s_service.go' % names)
		self.render_file_to('encode_service.go', file_path, names)

		file_path = os.path.join(dir_path, 'fill_%(resource)s_service.go' % names)
		self.render_file_to('fill_service.go', file_path, names)

	def render_file_to(self, template_name, target_dir_path, context):
		from jinja2 import Template

		with open('%s/eel/%s' % (TEMPLATE_FILE_DIR, template_name), 'rb') as f:
			template_content = f.read()

		template = Template(template_content.decode('utf-8'))
		content = template.render(context)

		print '> generate: ', target_dir_path
		with open(target_dir_path, 'wb') as f:
			print >> f, content.encode('utf-8')

	def copy_file(self, src, dst):
		print '> copy: %s -> %s' % (src, dst)
		if os.path.exists(dst):
			print '[WARN]: %s is already EXISTS!!!, DO NOT COPY.' % dst
			#shutil.copyfile(src, dst)
		else:
			shutil.copyfile(src, dst)

	def copy_files(self, context):
		print '\n>>>>>>>>>> generate business objects <<<<<<<<<<'
		src = '_generate/models/%(resource)s.go' % context
		dst = 'models/%(full_package)s/%(resource)s.go' % context
		self.copy_file(src, dst)

		src = '_generate/rest/%(resource)s.go' % context
		dst = 'rest/%(full_package)s/%(resource)s.go' % context
		self.copy_file(src, dst)

		src = '_generate/rest/%(plural_resource)s.go' % context
		dst = 'rest/%(full_package)s/%(plural_resource)s.go' % context
		self.copy_file(src, dst)

		src = '_generate/business/%(resource)s.go' % context
		dst = 'business/%(full_package)s/%(resource)s.go' % context
		self.copy_file(src, dst)

		src = '_generate/business/%(resource)s_repository.go' % context
		dst = 'business/%(full_package)s/%(resource)s_repository.go' % context
		self.copy_file(src, dst)

	def download_code_base(self, url, zipfile):
		total_bytes = 0
		with open(zipfile, 'wb') as handle:
			response = requests.get(url, stream=True)

			for block in response.iter_content(1024):
				total_bytes += len(block)
				print 'download......%sk' % round((total_bytes/1024.0), 2)
				handle.write(block)
				time.sleep(.001)

	def unzip_code_base_to(self, name):
		zfobj = zipfile.ZipFile(code_zip_path)
		zfobj.extractall()

		for dir in os.listdir('.'):
			if not os.path.isdir(dir):
				continue

			if not 'golang-service-resource-template' in dir:
				continue

			print 'rename %s to %s' % (dir, name)
			os.rename(dir, name)

		zfobj.close()
		os.remove(code_zip_path)

	def download_template_files(self):
		print 'download Golang template file'
		code_base_url = 'https://code.aliyun.com/clubxiaocheng/golang-service-resource-template/repository/archive.zip'
		self.download_code_base(code_base_url, code_zip_path)
		self.unzip_code_base_to(TEMPLATE_FILE_DIR)
	
	def handle(self, package, resource, **options):
		generated_dir = './_generate'
		if os.path.exists(generated_dir):
			shutil.rmtree(generated_dir)

		#self.download_template_files()
	
		names = self.get_names(resource)
		names['package'] = package
		names['full_package'] = package
		names['resource'] = resource
		self.generate_model_file(names)
		self.generate_resource_file(names)
		self.generate_business_file(names)

		# print '\n******************** Generate File ********************'
		# print 'file is generated under ./_generate dir, please copy to real dirs'
		# print 'Do you want to copy files now? (y/n): ',

		#if os.path.exists(TEMPLATE_FILE_DIR):
		#	print "remove %s" % TEMPLATE_FILE_DIR
		#	shutil.rmtree(TEMPLATE_FILE_DIR)

		input = 'n'#raw_input().strip()
		if input == 'Y' or input == 'y':
			self.copy_files(names)
			print '\n******************** Success ********************'
			print 'Modify `models/init.go`, `routers/router.go` to connect resource into system'
		else:
			print '\n******************** Success ********************'
			print 'NOT COPY FILE. Please copy files manually'
			print 'And modify `models/init.go`, `routers/router.go` to connect resource into system'
