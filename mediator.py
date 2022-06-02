import asyncio
import logging
import json_module
from server import *
from aiogram import Bot, Dispatcher, executor, types


with open("config.json", encoding='utf-8') as file:
        config = json.load(file)

token = config["token_mediator"]
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)
logger.setLevel(logging.DEBUG)

bot = Bot(token=token)
loop = asyncio.get_event_loop()
dp = Dispatcher(bot, loop=loop)


def send_message_to_observer(message):
    themes = json_module.open_json("themes")
    theme_id = str(server_get_theme_id_by_user(chat_id=message.chat['id']))
    result = f"[{message.date}] [{themes[theme_id]['name']}] @{message.chat['username']}: {message.text}"
    return result


@dp.message_handler(commands=['start'])
async def send_welcome(message: types.Message):
    themes = json_module.open_json("themes")
    theme_id = message.text
    theme_id = theme_id.split()[1]
    if theme_id in themes.keys():
        welcome_message = f'{themes[theme_id]["name"]}\nОписание: {themes[theme_id]["description"]}'
        try:
            server_create_user(message.chat['id'], theme_id, message.chat['first_name'], f"@{message.chat['username']}")
        except:
            server_update_user_theme(theme_id, message.chat['id'])

    else:
        welcome_message = "Темы не существует"
    await message.reply(welcome_message)


@dp.message_handler()
async def echo(message: types.Message, state=None):
    server_send_message(message.message_id, message.text, message.chat['id'], message.date, message.chat['username'], server_get_theme_id_by_user(message.chat['id']))
    msg_notif = send_message_to_observer(message)
    notif = json_module.open_json("notif")
    notif['value'] = msg_notif
    notif['chat_id'] = message.chat['id']
    json_module.update_json("notif", notif)

async def check_answers():
    while True:
        answer = json_module.open_json("answer")
        users_json = json_module.open_json("send_users")
        if answer['answer'] != "":
            await bot.send_message(answer['chat_id'], answer['answer'])
            answer['answer'] = ""
            answer['chat_id'] = ""
            json_module.update_json("answer", answer)
        if users_json['user_id_list'] != "":
            for i in range(len(users_json['user_id_list'])):
                await bot.send_message(users_json['user_id_list'][i], users_json['message'])
            users_json['user_id_list'] = ""
            users_json['message'] = ""
            json_module.update_json("send_users", users_json)

        await asyncio.sleep(0.1)

if __name__ == '__main__':
    dp.loop.create_task(check_answers())
    executor.start_polling(dp, skip_updates=True)