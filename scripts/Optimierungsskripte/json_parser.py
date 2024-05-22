import json

def loadJSON(stringJSON: str) -> dict:
  return json.loads(stringJSON)

def getGroup(dictJSON: dict) -> dict[str, list[str]]:
  temp = []
  for lecturer in dictJSON["lecturers"]:
    for lecture in lecturer["lectures"]:
      temp.append(lecture["name"])
  return {dictJSON["name"]:temp}

def getLectures(dictJSON: dict) -> dict[str, int]:
  temp = {}
  for lecturer in dictJSON["lecturers"]:
    for lecture in lecturer["lectures"]:
      temp[lecture["name"]] = lecture["amount"] # .append(lecture["name"])
  return temp
  

def getLecturers(dictJSON: dict) -> dict[str, list[str]]:
  temp = {}
  for lecturer in dictJSON["lecturers"]:
    temp[lecturer["name"]] = []
    for lecture in lecturer["lectures"]:
      temp[lecturer["name"]].append(lecture["name"]) # .append(lecture["name"])
  return temp

