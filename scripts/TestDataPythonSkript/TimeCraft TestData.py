import json
from bson import ObjectId  # Füge diese Zeile am Anfang deines Skripts hinzu
import os

def clear_terminal():
    # Je nach Betriebssystem wird der passende Befehl ausgeführt
    os.system('cls' if os.name == 'nt' else 'clear')


def get_user_data():
    print("Enter User Data:")
    email = input("Email: ")
    password = input("Password: ")
    is_verified_email = input("Is Verified Email (True/False): ").lower() == 'true'
    username = input("Username: ")
    permissions = input("Permissions (comma-separated ObjectIDs): ").split(',')
    ssotoken = input("SSOToken: ")

    user_data = {
        "email": email,
        "password": password,
        "isVerifiedEmail": is_verified_email,
        "username": username,
        "permissions": [int(perm) for perm in permissions],
        "SSOToken": ssotoken
    }

    return user_data

def get_room_config_data():
    print("Enter Room Config Data:")
    name = input("Name: ")
    room_nr = input("Room Number: ")
    capacity = int(input("Capacity: "))
    exam_capacity = int(input("Exam Capacity: "))
    blocks = input("Blocks (comma-separated ObjectIDs): ").split(',')
    specialisation = input("Specialisation (comma-separated ObjectIDs): ").split(',')

    room_config_data = {
        "name": name,
        "roomNr": room_nr,
        "capacity": capacity,
        "examCapacity": exam_capacity,
        "blocks": blocks,
        "specialisation": [int(spec) for spec in specialisation]
    }


    return room_config_data

def get_room_specialisation_data():
    print("Enter Room Specialisation Data:")
    name = input("Name: ")

    room_specialisation_data = {
        "name": name,
    }

    return room_specialisation_data


def get_lecture_data():
    print("Enter Lecture Data:")
    name = input("Name: ")
    contactHours = int(input("ContactHours: "))
    

    lecture_data = {
        "name": name,
        "contactHours": contactHours,
    }


    return lecture_data




def get_lecturer_data():
    print("Enter Lecturer Data:")
    firstName = input("First Name: ")
    lastName = input("Last Name: ")
    title = input("Title: ")
    canHold = input("canHold (comma-separated ObjectIDs): ").split(',')
    contactEmail = input("Contact Email: ")
    maxWeekHours = int(input("Max Week Hours: "))
    contactPhone = input("Contact Phone: ")

    lecturer_data = {
        "FirstName": firstName,
        "LastName": lastName,
        "Title": title,
        "CanHold": canHold,
        "ContactEmail": contactEmail,
        "MaxWeekHours": maxWeekHours,
        "ContactPhone": contactPhone,
    }

    return lecturer_data



def main():
    model_choice = input("Enter one of the following model:\nUser\nRoomconfig\nRoomSpecialisation\nLecture\nLecturer\n\n ").lower()
    clear_terminal


    if model_choice == "user":
        data = get_user_data()
        filename_key = data["username"]  # Hier verwende das gewünschte Schlüssel-Element
    elif model_choice == "roomconfig":
        data = get_room_config_data()
        filename_key = data["roomNr"]  # Hier verwende das gewünschte Schlüssel-Element
    elif model_choice == "roomspecialisation":
        data = get_room_specialisation_data()
        filename_key = data["name"]  # Hier verwende das gewünschte Schlüssel-Element
    elif model_choice == "lecture":
        data = get_lecture_data()
        filename_key = data["name"]  # Hier verwende das gewünschte Schlüssel-Element
    elif model_choice == "lecturer":
        data = get_lecturer_data()
        filename_key = data["LastName"]  # Hier verwende das gewünschte Schlüssel-Element
    else:
        print("Invalid model choice. Exiting.")
        return

    confirm = input(f"Confirm data entry for {model_choice.capitalize()} (Y/N): ").lower()
    if confirm == "y":
        with open(f"{filename_key}_data.json", "w") as json_file:
            json.dump(data, json_file)
        print(f"{model_choice.capitalize()} data saved to {filename_key}_data.json.")
    else:
        print("Data entry canceled.")


if __name__ == "__main__":
    main()
