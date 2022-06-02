import asyncio
import logging
from aiogram.dispatcher.filters import Text
import random
from server import *
import json_module
from aiogram import Bot, Dispatcher, executor, types
from aiogram.contrib.fsm_storage.memory import MemoryStorage
from aiogram.dispatcher.filters.state import State, StatesGroup
from aiogram.types import InlineKeyboardButton, InlineKeyboardMarkup


def create_theme(name, description, mails):
    themes = json_module.open_json("themes")
    theme_id = 0
    while theme_id in themes.keys() or theme_id == 0:
        theme_id = random.randint(10000000, 99999999)

    theme_link = "https://t.me/Mediator_system_bot?start=" + str(theme_id)
    themes[theme_id] = {'name': name, 'description': description}
    # mail
    #email_module.send_mail(name, description, mails, theme_link)
    # server
    server_create_theme(theme_id, name)

    json_module.update_json("themes", themes)
    return theme_link


def show_themes():
    message = ""
    themes = json_module.open_json("themes")
    i = 1
    themes_list = []
    for key in themes:
        message += f"{i}: {themes[key]['name']}\nСсылка: {'https://t.me/Mediator_system_bot?start='+key}\n"
        themes_list.append(key)
        i += 1
    return message, themes_list


def delete_theme(delete_key):
    themes = json_module.open_json("themes")
    themes.pop(delete_key)
    # server
    server_delete_theme(delete_key)
    json_module.update_json("themes", themes)

# messages

def show_all_messages(messages):
    result = ''
    themes = json_module.open_json("themes")
    for i in range(len(messages)):
        result += f"[{messages[i]['time'][:19:]}] [{themes[str(messages[i]['theme_id'])]['name']}] @{messages[i]['ref']}: {messages[i]['message']}\n"
    return result

async def check_notifications():
    while True:
        notif = json_module.open_json("notif")
        if notif['value'] != "":
            admins = json_module.open_json("admins")
            for id in admins:
                inline_answer_mark = InlineKeyboardMarkup(row_width=4)
                inline_answer_mark.add(InlineKeyboardButton('Ответить', callback_data=f"answer_{notif['chat_id']}"))
                await bot.send_message(id, notif['value'], reply_markup=inline_answer_mark)
            notif['value'] = ""
            notif['chat_id'] = ""
            json_module.update_json("notif", notif)
        await asyncio.sleep(0.1)


class answerFSM(StatesGroup):
    answer = State()

class classFSM(StatesGroup):
    name = State()
    description = State()
    mailing_list = State()

class send_by_themesFSM(StatesGroup):
    message = State()

storage = MemoryStorage()

with open("config.json", encoding='utf-8') as file:
    config = json.load(file)

themes = json_module.open_json("themes")

API_TOKEN = config['token_observer']
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)
logger.setLevel(logging.DEBUG)

bot = Bot(token=API_TOKEN)
loop = asyncio.get_event_loop()
dp = Dispatcher(bot, storage=storage, loop=loop)

@dp.message_handler(commands=['start'])
async def send_welcome(message: types.Message):
    admins = json_module.open_json("admins")
    if message.chat['id'] not in admins:
        admins.append(message.chat['id'])
    json_module.update_json("admins", admins)
    keyboard_markup = types.ReplyKeyboardMarkup(row_width=3, resize_keyboard=True)
    more_btns_text = (
        "Создать тему",
        "Показать темы",
        "Показать все сообщения",
        "Показать сообщения по теме",
        "Рассылка по теме"
    )
    keyboard_markup.add(*(types.KeyboardButton(text) for text in more_btns_text))
    await message.reply("Используйте меню для работы с ботом", reply_markup=keyboard_markup)


@dp.callback_query_handler(Text(startswith="numdel_"))
async def process_callback_button(callback_query: types.CallbackQuery):
    reply_text, themes_list = show_themes()
    delete_key = None
    action = callback_query.data.split("_")[1]
    for i in range(1, len(themes_list) + 1):
        if i == int(action):
            delete_key = themes_list[i-1]
            delete_theme(delete_key)
    await bot.answer_callback_query(callback_query.id)
    await bot.send_message(callback_query.from_user.id, f'Удалено {action}')


@dp.callback_query_handler(Text(startswith="numshow_"))
async def process_callback_button(callback_query: types.CallbackQuery):
    reply_text, themes_list = show_themes()
    check_key = None
    action = callback_query.data.split("_")[1]
    for i in range(1, len(themes_list) + 1):
        if i == int(action):
            check_key = themes_list[i - 1]
    await bot.send_message(callback_query.from_user.id, f"{show_all_messages(server_get_message_by_theme(check_key))}")


@dp.callback_query_handler(Text(startswith="numsendtheme_"))
async def process_callback_button(callback_query: types.CallbackQuery):
    reply_text, themes_list = show_themes()
    check_key = None
    action = callback_query.data.split("_")[1]
    for i in range(1, len(themes_list) + 1):
        if i == int(action):
            check_key = themes_list[i - 1]
    await send_by_themesFSM.message.set()
    state = Dispatcher.get_current().current_state()
    async with state.proxy() as data:
        data['theme_id'] = check_key
    await bot.send_message(callback_query.from_user.id, f'Введите сообщение для рассылки:')



@dp.callback_query_handler(Text(startswith="answer_"))
async def process_callback_button(callback_query: types.CallbackQuery):
    chat_id = callback_query.data.split("_")[1]
    await bot.answer_callback_query(callback_query.id)
    await answerFSM.answer.set()
    state = Dispatcher.get_current().current_state()
    async with state.proxy() as data:
        data['chat_id'] = chat_id
    await bot.send_message(callback_query.from_user.id, f'Введите ответ:')


@dp.message_handler(state=send_by_themesFSM)
async def load_answer(message: types.Message, state: send_by_themesFSM):
    async with state.proxy() as data:
        data['message'] = message.text
        users_json = json_module.open_json("send_users")
        users_json['user_id_list'] = server_get_users_id_by_theme(data['theme_id'])
        users_json['message'] = data['message']
        json_module.update_json("send_users", users_json)
    await message.reply("Рассылка выполнена")
    await state.finish()

@dp.message_handler(state=answerFSM.answer)
async def load_answer(message: types.Message, state: answerFSM):
    async with state.proxy() as data:
        data['answer'] = message.text
        answer_json = json_module.open_json("answer")
        answer_json['answer'] = data['answer']
        answer_json['chat_id'] = data['chat_id']
        json_module.update_json("answer", answer_json)

    await message.reply("Ответ отправлен")
    await state.finish()


@dp.message_handler()
async def echo(message: types.Message, state=None):
    button_text = message.text
    reply_text = ''
    reply_markup = None
    if button_text == 'Создать тему':
        await classFSM.name.set()
        await message.reply('Введите название темы:')

    elif button_text == 'Показать темы':
        reply_text, themes_list = show_themes()
        inline_kb = InlineKeyboardMarkup(row_width=4)
        btn = []
        for i in range(1, len(themes_list) + 1):
            btn.append(InlineKeyboardButton(f'Удалить {i}', callback_data=f'numdel_{i}'))
        inline_kb.add(*btn)
        reply_markup = inline_kb

    elif button_text == 'Показать все сообщения':
        reply_text = show_all_messages(server_get_all_message())
    elif button_text == 'Показать сообщения по теме':
        reply_text, themes_list = show_themes()
        inline_msg_kb = InlineKeyboardMarkup(row_width=4)
        btn = []
        for i in range(1, len(themes_list) + 1):
            btn.append(InlineKeyboardButton(f'Посмотреть {i}', callback_data=f'numshow_{i}'))
        inline_msg_kb.add(*btn)
        reply_markup = inline_msg_kb
    elif "Рассылка по теме":
        reply_text, themes_list = show_themes()
        inline_sendthm_kb = InlineKeyboardMarkup(row_width=4)
        btn = []
        for i in range(1, len(themes_list) + 1):
            btn.append(InlineKeyboardButton(f'Разослать {i}', callback_data=f'numsendtheme_{i}'))
        inline_sendthm_kb.add(*btn)
        reply_markup = inline_sendthm_kb

    await message.reply(reply_text, reply_markup=reply_markup)

@dp.message_handler(state=classFSM.name)
async def load_name(message: types.Message, state: classFSM):
    async with state.proxy() as data:
        data['name'] = message.text
    await classFSM.next()
    await message.reply("Введите описание темы:")


@dp.message_handler(state=classFSM.description)
async def load_description(message: types.Message, state: classFSM):
    async with state.proxy() as data:
        themes = json_module.open_json("themes")
        data['description'] = message.text
    # TODO send to bd
    await message.reply("Введите email адреса для рассылки(через пробел):")
    await classFSM.next()

@dp.message_handler(state=classFSM.mailing_list)
async def load_mails(message: types.Message, state: classFSM):
    async with state.proxy() as data:
        data['mails'] = message.text
    await message.reply(f"Ссылка на новую тему: {create_theme(data['name'], data['description'], data['mails'])}")
    await state.finish()



if __name__ == '__main__':
    dp.loop.create_task(check_notifications())
    executor.start_polling(dp, skip_updates=True)
