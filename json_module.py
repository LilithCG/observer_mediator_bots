import json

def open_json(jsonfile):
    with open(f"config/{jsonfile}.json", encoding='utf-8') as file:
        answer = json.load(file)
    return answer

def update_json(jsonfile, data):
    with open(f"config/{jsonfile}.json", "w") as write_file:
        json.dump(data, write_file)