# -*- coding: utf-8 -*-

import sys
import json
import os
import json

import requests

__all__ = ('ServiceClient')

SERVICE2PORT = {
	'skep': '8301',
	'peanut': '8004',
	'coral': '8321',
	'plutus': '8006',
	'eel': '3131'
}

class FakeResponse(object):
	def __init__(self, method='get', url=None):
		self.body = None
		self.status = None
		self.headers = None
		self.cookies = None
		self.method = method
		self.url = url

	def __str__(self):
		buf = []
		buf.append('\n===== start response =====')
		buf.append('%s: %s' % (self.method, self.url))
		buf.append('*** body ***')
		if self.body and 'queries' in self.body:
			del self.body['queries']
		traceback = None
		if self.body and 'Traceback' in self.body.get('innerErrMsg'):
			# print str(self.body['innerErrMsg'].encode('utf-8'))
			# raw_input()
			traceback = self.body['innerErrMsg'].encode('utf-8')
			self.body['innerErrMsg'] = 'see traceback underneath...'
		buf.append(json.dumps(self.body, indent=2).decode('utf-8'))

		if traceback:
			lines = traceback.decode('utf-8').split('\n')
			count = len(lines)
			i = 0
			traceback_buf = []
			suspension_points = []
			while i < count:
				line = lines[i]
				if 'File ' in line:
					if 'site-packages' in line or 'features/util' in line:
						suspension_points.append('..')
						i += 2
						continue
					else:
						if len(suspension_points) != 0:
							traceback_buf.append('  ' + ''.join(suspension_points).decode('utf-8'))
							suspension_points = []
						traceback_buf.append(line)
						i += 1
						traceback_buf.append(lines[i])
						i += 1
						continue
				else:
					traceback_buf.append(line)
					i += 1
			buf.append('** exception traceback **')
			buf.append(u'\n'.join(traceback_buf))
		buf.append('*** http status ***')
		#buf.append(self.status)
		buf.append('*** headers ***')
		buf.append('===== end response =====\n')
		result = u'\n'.join(buf)
		return result.encode('utf-8')

	@property
	def status_code(self):
		return self.status

	@property
	def data(self):
		return self.body['data']

	@property
	def is_success(self):
		return self.body['code'] == 200

class ServiceClient(object):
	def __get_url(self, resource):
		pos = resource.find(':')
		if pos == -1:
			raise ValueError('invalid remote resource: %s' % resource)

		service = resource[:pos]
		resource = resource[pos+1:]

		pos = resource.rfind('.')
		app = resource[:pos]
		resource = resource[pos+1:]

		host = '127.0.0.1'
		return "http://{}:{}/{}/{}/?_v=1".format(host, SERVICE2PORT[service], app, resource)

	def get(self, resource, data={}, jwt_token=''):
		url = self.__get_url(resource)

		headers = {
			'AUTHORIZATION': jwt_token
		}
		response = requests.get(url, headers=headers, params=data)
		
		fake_response = FakeResponse('GET', url)
		fake_response.body = json.loads(response.content)
		return fake_response

	def put(self, resource, data={}, jwt_token=''):
		url = self.__get_url(resource)
		url = '%s&_method=put' % url

		headers = {
			'AUTHORIZATION': jwt_token
		}
		response = requests.post(url, data, headers=headers)

		fake_response = FakeResponse('PUT', url)
		fake_response.body = json.loads(response.content)
		return fake_response

	def post(self, resource, data={}, jwt_token=''):
		url = self.__get_url(resource)

		headers = {
			'AUTHORIZATION': jwt_token
		}
		response = requests.post(url, data, headers=headers)

		fake_response = FakeResponse('POST', url)
		fake_response.body = json.loads(response.content)
		return fake_response

	def delete(self, resource, data={}, jwt_token=''):
		url = self.__get_url(resource)
		url = '%s&_method=delete' % url

		headers = {
			'AUTHORIZATION': jwt_token
		}
		response = requests.post(url, data, headers=headers)

		fake_response = FakeResponse('DELETE', url)
		fake_response.body = json.loads(response.content)
		return fake_response