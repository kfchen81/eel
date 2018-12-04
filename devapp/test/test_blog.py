# -*- coding: utf-8 -*-

from service_client import ServiceClient

JWT = 'eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJjb2RlIjogIjMwMDAwMDMxOCIsICJ1aWQiOiAxMTA0LCAicGlkIjogMiwgIm1pZCI6IDEsICJleHAiOiAiMjAxOS0xMS0yNyAyMDoxMyIsICJ2IjogMTgwNzE5LCAiaWF0IjogIjIwMTgtMTEtMjcgMjA6MTMiLCAidHlwZSI6IDIsICJwbCI6IDU2fQ==.yZZOuVOsQY5w7URSaQaOyfkP8KVkOcNE+s5uP0lU6mE='

client = ServiceClient()

# print "========== PUT =========="
# resp = client.put("eel:blog.blog", {
# 	"title": "blog title 1",
# 	"content": "blog content 1"
# }, jwt_token=JWT)
# print resp.body['code']

print "========== GET FIRST BLOG =========="
resp = client.get("eel:blog.blog", {
}, jwt_token=JWT)
print resp.body['code']
blog = resp.body['data']

print "========== UPDATE COUNTER =========="
resp = client.post("eel:blog.counter", {
	"id": blog['id'],
	"delta": 5
}, jwt_token=JWT)
print resp.body['code']
#
# print "========== GET =========="
# resp = client.get("eel:blog.blogs", {
# 	'_p_from': 99999999,
# 	'_p_count': 2
# }, jwt_token=JWT)
# first_blog = resp.body['data']['blogs'][0]
# print first_blog['id'], ' ', first_blog['title'], ' ', first_blog['content']
#
# blog_id = resp.body['data']['blogs'][0]['id']
# print "========== PUT =========="
# resp = client.post("eel:blog.blog", {
# 	'id': blog_id,
# 	'title': 'updated title',
# 	'content': 'updated content'
# }, jwt_token=JWT)
# print resp.body['code']
# resp = client.get("eel:blog.blogs", {
# 	'_p_from': 99999999,
# 	'_p_count': 2
# }, jwt_token=JWT)
# first_blog = resp.body['data']['blogs'][0]
# print first_blog['id'], ' ', first_blog['title'], ' ', first_blog['content']
#
# print "========== DELETE =========="
# resp = client.delete("eel:blog.blog", {
# 	'id': blog_id
# }, jwt_token=JWT)
# print resp.body['code']
# resp = client.get("eel:blog.blogs", {
# 	'_p_from': 99999999,
# 	'_p_count': 2
# }, jwt_token=JWT)
# first_blog = resp.body['data']['blogs'][0]
# print first_blog['id'], ' ', first_blog['title'], ' ', first_blog['content']
