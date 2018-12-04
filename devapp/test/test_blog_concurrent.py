# -*- coding: utf-8 -*-
#__author__ = 'robert'
from datetime import datetime
import threading
import time
import random
import requests
from service_client import ServiceClient


THREAD_COUTN = 1
ITERATE_COUNT = 1
DUMP_STEPS = 10
THREAD2ACTIONS = {}
THREAD2COUNT = {}

JWT = 'eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJjb2RlIjogIjMwMDAwMDMxOCIsICJ1aWQiOiAxMTA0LCAicGlkIjogMiwgIm1pZCI6IDEsICJleHAiOiAiMjAxOS0xMS0yNyAyMDoxMyIsICJ2IjogMTgwNzE5LCAiaWF0IjogIjIwMTgtMTEtMjcgMjA6MTMiLCAidHlwZSI6IDIsICJwbCI6IDU2fQ==.yZZOuVOsQY5w7URSaQaOyfkP8KVkOcNE+s5uP0lU6mE='


class ResultAnalyzer(object):
	def __init__(self):
		pass

	def analyze(self):
		action2cumcount = {
			"get": 0,
			"put": 0,
			"post": 0,
			"delete": 0
		}
		for thread_id, action2count in THREAD2ACTIONS.items():
			for action, count in action2count.items():
				action2cumcount[action] = action2cumcount[action] + count

		print action2cumcount

		#counter
		total_count = 0
		for thread_id, count in THREAD2COUNT.items():
			total_count += count
		print 'total counter: ', total_count

class ConsumeAdThread(threading.Thread):
	def __init__(self, id):
		threading.Thread.__init__(self)
		self.id = id
		self.client = ServiceClient()

	def get(self):
		resp = self.client.get("eel:blog.blogs", {
			'_p_from': 99999999,
			'_p_count': 2
		}, jwt_token=JWT)
		if resp.body['code'] != 200:
			print '[Error] GET - %d' % self.id
		else:
			global THREAD2ACTIONS
			THREAD2ACTIONS[self.id]["get"] = THREAD2ACTIONS[self.id]["get"]+1

	def create(self, index):
		resp = self.client.put("eel:blog.blog", {
			"title": "blog title %d-%d" % (self.id, index),
			"content": "blog content"
		}, jwt_token=JWT)
		if resp.body['code'] != 200:
			print '[Error] GET - %d' % self.id
		else:
			global THREAD2ACTIONS
			THREAD2ACTIONS[self.id]["put"] = THREAD2ACTIONS[self.id]["put"]+1

	def update(self):
		resp = self.client.get("eel:blog.blogs", {
			'_p_from': 99999999,
			'_p_count': 2
		}, jwt_token=JWT)
		if len(resp.body['data']['blogs']) == 0:
			return

		resp = self.client.get("eel:blog.blog", {
		}, jwt_token=JWT)
		if resp.body['code'] != 200:
			print '[Error] GET - %d' % self.id

		#only update when have at least one blog
		first_blog = resp.body['data']
		delta = random.randint(1, 100)
		resp = self.client.post("eel:blog.counter", {
			"id": first_blog["id"],
			"delta": delta
		}, jwt_token=JWT)
		if resp.body['code'] != 200:
			print '[Error] GET - %d' % self.id
		else:
			global THREAD2ACTIONS
			THREAD2ACTIONS[self.id]["post"] = THREAD2ACTIONS[self.id]["post"]+1

			global THREAD2COUNT
			THREAD2COUNT[self.id] = THREAD2COUNT[self.id] + delta

	def delete(self):
		resp = self.client.get("eel:blog.blogs", {
			'_p_from': 99999999,
			'_p_count': 2
		}, jwt_token=JWT)
		if len(resp.body['data']['blogs']) == 0:
			return

		resp = self.client.get("eel:blog.blog", {
		}, jwt_token=JWT)
		if resp.body['code'] != 200:
			print '[Error] GET - %d' % self.id

		#only update when have at least one blog
		first_blog = resp.body['data']
		delta = random.randint(-100, 0)
		resp = self.client.post("eel:blog.counter", {
			"id": first_blog["id"],
			"delta": delta
		}, jwt_token=JWT)
		if resp.body['code'] != 200:
			print '[Error] GET - %d' % self.id
		else:
			global THREAD2ACTIONS
			THREAD2ACTIONS[self.id]["delete"] = THREAD2ACTIONS[self.id]["delete"]+1

			global THREAD2COUNT
			THREAD2COUNT[self.id] = THREAD2COUNT[self.id] + delta

	def run(self):
		random.seed(time.time())
		for i in range(0, ITERATE_COUNT):
			index = random.randint(1, 100)
			if index > 90:
				self.delete()
			elif index > 50:
				self.update()
			elif index > 30:
				self.get()
			else:
				self.create(i)
			if i % DUMP_STEPS == 0:
				print '[%d] - %s' % (self.id, i)
			# sleep_delta = random.randint(1, 7)
			# time.sleep(sleep_delta/10.0)

class Command(object):
	help = "python manage.py test_ad_concurrenty [thread_count] [iterate_count]"
	args = ''

	def __init__(self):
		self.threads = []

	def handle(self, thread_count, iterate_count):
		global THREAD_COUTN
		THREAD_COUTN = int(thread_count)

		global ITERATE_COUNT
		ITERATE_COUNT = int(iterate_count)

		global THREAD2ACTIONS
		start = time.time()
		for i in range(0, THREAD_COUTN):
			THREAD2ACTIONS[i] = {
				"put": 0,
				"post": 0,
				"get": 0,
				"delete": 0
			}
			THREAD2COUNT[i] = 0
			thread = ConsumeAdThread(i)
			self.threads.append(thread)
			thread.start()

		for t in self.threads:
			t.join()
		end = time.time()
		print 'threads finished: %ss' % int((end - start))

		try:
			result_analyzer = ResultAnalyzer()
			result_analyzer.analyze()
		except:
			import sys, traceback
			type, value, tb = sys.exc_info()
			print type
			print value
			traceback.print_tb(tb)

command = Command()
command.handle(100, 100)