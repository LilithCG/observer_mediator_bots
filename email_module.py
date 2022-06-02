import json
import smtplib
from email.mime.multipart import MIMEMultipart
from email.mime.text import MIMEText

with open("config.json", encoding='utf-8') as file:
    config = json.load(file)


def send_mail(topic, description, mails, link):
    mails = mails.split()
    address_from = config['email_from']
    password = config['email_password']

    server = smtplib.SMTP('smtp.gmail.com: 587')
    server.starttls()
    server.login(address_from, password)

    for mail in range(len(mails)):
        try:
            msg = MIMEMultipart()
            msg['From'] = address_from
            msg['To'] = mails[mail]
            msg['Subject'] = topic
            msg.attach(MIMEText(description + "\nСсылка для связи: " + link, 'plain'))
            server.sendmail(msg['From'], msg['To'], msg.as_string())

        except:
            print("SEND_MAIL_ERROR")


    server.quit()
