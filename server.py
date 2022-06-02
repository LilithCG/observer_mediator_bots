import json

import pymysql
import requests


def server_create_theme(theme_id, name):
    requests.post(f'http://127.0.0.1:5000/api/createtheme?theme_id={theme_id}&name={name}')


def server_delete_theme(theme_id):
    requests.delete(f'http://127.0.0.1:5000/api/deletetheme?id={theme_id}')


def server_create_user(chat_id, theme_id, name, ref):
    requests.post(f'http://127.0.0.1:5000/api/createuser?chat_id={chat_id}&name={name}&reference={ref}&theme_id={theme_id}')


def server_send_message(id, message, chat_id, date, ref, theme_id):
    requests.post(f'http://127.0.0.1:5000/api/createmsg?id={id}&chat_id={chat_id}&theme_id={theme_id}&message={message}&ref={ref}')

def server_get_all_message():
    messages = requests.get(f'http://127.0.0.1:5000/api/getall').content
    result = messages
    result = result.decode('utf8').replace("'", '"')
    result = json.loads(result)
    return result


def server_update_user_theme(theme_id, chat_id):
    requests.put(f'http://127.0.0.1:5000/api/update?theme_id={theme_id}&chat_id={chat_id}')


def server_get_message_by_theme(theme_id):
    messages = requests.get(f'http://127.0.0.1:5000/api/get?id={theme_id}').content
    result = messages
    result = result.decode('utf8').replace("'", '"')
    result = json.loads(result)
    return result

def server_get_theme_id_by_user(chat_id):
    messages = requests.get(f'http://127.0.0.1:5000/api/get_users_chat?chat_id={chat_id}').content
    result = messages
    result = result.decode('utf8').replace("'", '"')
    result = json.loads(result)
    result = result[0]['theme_id']
    return result

def server_get_users_id_by_theme(theme_id):
    messages = requests.get(f'http://127.0.0.1:5000/api/get_users_theme?theme_id={theme_id}').content
    result = messages
    result = result.decode('utf8').replace("'", '"')
    result = json.loads(result)
    res_list = []
    for i in range(len(result)):
        res_list.append(result[i]['chat_id'])
    return res_list