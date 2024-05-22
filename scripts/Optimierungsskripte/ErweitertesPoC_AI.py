from ortools.sat.python import cp_model
import math
from Room import Room
import json
import json_parser
from tabulate import tabulate

# define functions

# hard constraints
def create_weeks(num_weeks: int) -> list:
    """
    Define the number of weeks and create them.
    """
    weeks = []
    for week in range(num_weeks):
        weeks.append(f"KW{week+1}")
    return weeks

def create_timeslots(num_timeslots: int) -> list:
    """
    Create the timeslots.
    """
    return range(num_timeslots)

def create_schedule_vars(lectures: dict, weeks: list, timeslots: list, rooms: list, lecturers: dict, student_groups: dict, model: cp_model.CpModel) -> dict:
    """
    Create the scheduling variables.
    We create a dictionary where the key is a tuple (lecture, week, timeslot, room, lecturer, group)
    and the value is a boolean variable that indicates whether the lecture is scheduled at the specified timeslot and room by the specified lecturer for the specified group.
    """
    schedule = {}
    for lecture in lectures:
        for week in weeks:
            for timeslot in timeslots:
                for room in rooms:
                    for lecturer in lecturers:
                        for group in student_groups:
                            # Only create variables for lecturers skilled in the lecture and groups requiring the lecture
                            if (
                                lecture in lecturers[lecturer]
                                and lecture in student_groups[group]
                            ):
                                # We create multiple variables for the same lecture if its frequency is more than 1.
                                for _ in range(lectures[lecture]):
                                    schedule[
                                        (lecture, week, timeslot, room, lecturer, group)
                                    ] = model.NewBoolVar(
                                        "schedule_%s_%s_%i_%s_%s_%s"
                                        % (lecture, week, timeslot, room, lecturer, group)
                                    )
    return schedule

def add_each_lecture_is_assigned_hard_constraint(lectures: dict, weeks: list, timeslots: list, rooms: list, lecturers: dict, student_groups: dict, model: cp_model.CpModel) -> None:
    """
    (Hard Constraint) Each lecture is assigned to exactly the amount of timeslots, rooms and lecturers as its frequency for each group.
    """
    for lecture in lectures:
        for group in student_groups:
            if lecture in student_groups[group]:
                model.Add(
                    sum(
                        schedule[(lecture, week, timeslot, room, lecturer, group)]
                        for week in weeks
                        for timeslot in timeslots
                        for room in rooms
                        for lecturer in lecturers
                        if lecture in lecturers[lecturer]
                    )
                    == lectures[lecture]
                )

def add_lectures_can_not_be_in_same_timeslot_and_room_in_the_same_week_hard_constraint(lectures: dict, weeks: list, timeslots: list, rooms: list, lecturers: dict, student_groups: dict, model: cp_model.CpModel) -> None:
    """
    (Hard Constraint) Lectures can't be in the same timeslot and room on the same week.
    """
    for week in weeks:
        for timeslot in timeslots:
            for room in rooms:
                model.Add(
                    sum(
                        schedule[(lecture, week, timeslot, room, lecturer, group)]
                        for lecture in lectures
                        for lecturer in lecturers
                        for group in student_groups
                        if lecture in lecturers[lecturer]
                        and lecture in student_groups[group]
                        and "online" not in room.name
                    )
                    <= 1
                )

def add_lecturer_can_not_be_assigned_twice_hard_constraint(lectures: dict, weeks: list, timeslots: list, rooms: list, lecturers: dict, student_groups: dict, model: cp_model.CpModel) -> None:
    """
    (Hard Constraint) A lecturer can't be assigned twice.
    """
    for week in weeks:
        for timeslot in timeslots:
            for lecturer in lecturers:
                model.Add(
                    sum(
                        schedule[(lecture, week, timeslot, room, lecturer, group)]
                        for lecture in lectures
                        for room in rooms
                        for group in student_groups
                        if lecture in lecturers[lecturer]
                        and lecture in student_groups[group]
                    )
                    <= 1
                )

def add_student_group_can_not_be_assigned_twice_hard_constraint(lectures: dict, weeks: list, timeslots: list, rooms: list, lecturers: dict, student_groups: dict, model: cp_model.CpModel) -> None:
    """
    (Hard Constraint) A group can't be assigned twice.
    """
    for week in weeks:
        for timeslot in timeslots:
            for group in student_groups:
                model.Add(
                    sum(
                        schedule[(lecture, week, timeslot, room, lecturer, group)]
                        for lecture in lectures
                        for room in rooms
                        for lecturer in lecturers
                        if lecture in lecturers[lecturer]
                        and lecture in student_groups[group]
                    )
                    <= 1
                )

def add_lecturers_can_veto_days_hard_constraint(lectures: dict, weeks: list, timeslots: list, rooms: list, lecturers: dict, student_groups: dict, model: cp_model.CpModel) -> None:
    """ 
    (Hard Constraint) A Lecturer won't be assigned if they veto.
    """
    for veto_lecturer in vetos:
        for veto_week in vetos[veto_lecturer]:
            for slot in vetos[veto_lecturer][veto_week]:
                model.Add(
                    sum(
                        schedule[(lecture, veto_week, slot, room, veto_lecturer, group)]
                        for lecture in lectures
                        for room in rooms
                        for group in student_groups
                        if lecture in lecturers[veto_lecturer]
                        and lecture in student_groups[group]
                    )
                    == 0
                )

def add_online_and_on_site_lecture_must_not_be_on_same_day_per_group_hard_constraint(lectures: dict, weeks: list, timeslots: list, rooms: list, lecturers: dict, student_groups: dict, model: cp_model.CpModel) -> None:
    """ 
    (Hard Constraint)online and on-site lecture must not be on the same day (per group)
    """
    for week in weeks:
        for timeslot in timeslots:
            if timeslot % 3 == 1:
                for group in student_groups:
                    online_lecture = sum(
                        schedule[(lecture, week, timeslot, raum_online, lecturer, group)]
                        for lecture in lectures
                        for lecturer in lecturers
                        if lecture in lecturers[lecturer]
                        and lecture in student_groups[group]
                    )
                    
                    adjacent_on_site_lecture = sum(
                        schedule[(lecture, week, timeslot + 1, room, lecturer, group)]
                        for lecture in lectures
                        for lecturer in lecturers
                        for room in rooms
                        if lecture in lecturers[lecturer]
                        and lecture in student_groups[group]
                        and room != raum_online
                    ) + sum(
                        schedule[(lecture, week, timeslot - 1, room, lecturer, group)]
                        for lecture in lectures
                        for lecturer in lecturers
                        for room in rooms
                        if lecture in lecturers[lecturer]
                        and lecture in student_groups[group]
                        and room != raum_online
                    )
                    
                    # Create a boolean variable for the adjacent on-site lecture condition
                    has_adjacent_on_site_lecture = model.NewBoolVar(f'has_adjacent_on_site_lecture_{week}_{timeslot}_{group}')
                    
                    model.Add(adjacent_on_site_lecture > 0).OnlyEnforceIf(has_adjacent_on_site_lecture)
                    model.Add(adjacent_on_site_lecture == 0).OnlyEnforceIf(has_adjacent_on_site_lecture.Not())
                    
                    model.Add(online_lecture == 0).OnlyEnforceIf(has_adjacent_on_site_lecture)

def add_max_30_percent_online_per_group_hard_constraint(lectures: dict, weeks: list, timeslots: list, rooms: list, lecturers: dict, student_groups: dict, model: cp_model.CpModel) -> None:
    """
    (Hard Constraint) max 30% online per group
    """
    for group in student_groups:
        total_lectures = sum(lectures[lecture] for lecture in lectures if lecture in student_groups[group])
        total_online_lectures = math.floor(total_lectures * 0.3)
        model.Add(sum(schedule.get((lecture, week, timeslot, raum_online, lecturer, group))
                    for lecture in lectures
                    for week in weeks
                    for timeslot in timeslots
                    for lecturer in lecturers
                    if lecture in lecturers[lecturer]) <= total_online_lectures)

# soft constraints
def add_minimize_first_block_per_day_soft_constraint(penalty_weight: int, lectures: dict, weeks: list, timeslots: list, rooms: list, lecturers: dict, student_groups: dict, model: cp_model.CpModel) -> cp_model.IntegralT | cp_model.LinearExpr:
    """ 
    (Soft Constraint) Minimize the usage of the first block per day
    """
    penalty_first_block = model.NewIntVar(0, 100, "penalty_first_block")
    first_block_usage = sum(
        schedule.get((lecture, week, timeslot, room, lecturer, group), 0)
        for lecture in lectures
        for week in weeks
        for room in rooms
        for lecturer in lecturers
        for group in student_groups
        for timeslot in timeslots
        if timeslot % 3 == 0
    )
    model.Add(penalty_first_block == first_block_usage)

    objective_function_with_penalty_first = penalty_first_block * penalty_weight
    return objective_function_with_penalty_first

def add_minimize_third_block_per_day_soft_constraint(penalty_weight: int, lectures: dict, weeks: list, timeslots: list, rooms: list, lecturers: dict, student_groups: dict, model: cp_model.CpModel) -> cp_model.IntegralT | cp_model.LinearExpr:
    """
    (Soft Constraint) Minimize the usage of the third block per day
    """
    penalty_third_block = model.NewIntVar(0, 100, "penalty_third_block")
    third_block_usage = sum(
        schedule.get((lecture, week, timeslot, room, lecturer, group), 0)
        for lecture in lectures
        for week in weeks
        for room in rooms
        for lecturer in lecturers
        for group in student_groups
        for timeslot in timeslots
        if timeslot % 3 == 2
    )
    model.Add(penalty_third_block == third_block_usage)

    objective_function_with_penalty_third = penalty_third_block * penalty_weight

    return objective_function_with_penalty_third

def add_minimize_saturday_lectures_soft_constraint(penalty_weight: int, lectures: dict, weeks: list, timeslots: list, rooms: list, lecturers: dict, student_groups: dict, model: cp_model.CpModel) -> cp_model.IntegralT | cp_model.LinearExpr:
    """
    (Soft Constraint) Minimize the usage of saturdays meaning blocks 15,16,17
    """ 
    penalty_saturday_block = model.NewIntVar(0, 100, "penalty_saturday_block")
    saturday_block_usage = sum(
        schedule.get((lecture, week, timeslot, room, lecturer, group), 0)
        for lecture in lectures
        for week in weeks
        for room in rooms
        for lecturer in lecturers
        for group in student_groups
        for timeslot in timeslots
        if timeslot > 14
    )
    model.Add(penalty_saturday_block == saturday_block_usage)
    objective_function_with_penalty_saturday = (
        penalty_saturday_block * penalty_weight
    )
    return objective_function_with_penalty_saturday

def add_minimize_rooms_with_high_ranking_soft_constraint(penalty_weight: int, lectures: dict, weeks: list, timeslots: list, rooms: list, lecturers: dict, student_groups: dict, model: cp_model.CpModel) -> cp_model.IntegralT | cp_model.LinearExpr:
    """ 
    (Soft Constraint) Minimize the usage of rooms with a high ranking number
    """
    penalty_high_ranking = model.NewIntVar(0, 100, "penalty_high_ranking")
    penalty_high_ranking_usage = sum(
        schedule.get((lecture, week, timeslot, room, lecturer, group), 0)
        for week in weeks
        for lecture in lectures
        for room in rooms
        for lecturer in lecturers
        for group in student_groups
        for timeslot in timeslots
        if room.ranking > 1
    )
    model.Add(penalty_high_ranking == penalty_high_ranking_usage)
    objective_function_with_penalty_high_ranking = (penalty_high_ranking * penalty_weight)

    return objective_function_with_penalty_high_ranking

def add_minimize_lectures_in_last_three_weeks_soft_constraint(penalty_weight: int, exam_weeks: int, lectures: dict, weeks: list, timeslots: list, rooms: list, lecturers: dict, student_groups: dict, model: cp_model.CpModel) -> cp_model.IntegralT | cp_model.LinearExpr:
    """
    (Soft Constraint) Minimize lectures during last Exam Weeks (time of exams)
    """
    penalty_last_three_weeks = model.NewIntVar(0, 100, "penalty_last_three_weeks_block")
    penalty_last_three_weeks_usage = sum(
        schedule.get((lecture, week, timeslot, room, lecturer, group), 0)
        for week in weeks[-exam_weeks:]
        for lecture in lectures
        for room in rooms
        for lecturer in lecturers
        for group in student_groups
        for timeslot in timeslots
    )
    model.Add(penalty_last_three_weeks == penalty_last_three_weeks_usage)
    objective_function_with_penalty_last_three_weeks = (
        penalty_last_three_weeks * penalty_weight
    )
    return objective_function_with_penalty_last_three_weeks

def add_minimize_online_lectures_soft_constraint(penalty_weight: int, lectures: dict, weeks: list, timeslots: list, room_online: Room, lecturers: dict, student_groups: dict, model: cp_model.CpModel) -> cp_model.IntegralT | cp_model.LinearExpr:
    """
    (Soft Constraint) Minimize the online lectures
    """
    penalty_online_lecture = model.NewIntVar(0, 100, "penalty_online_lecture")
    penalty_online_lecture_usage = sum(
        schedule.get((lecture, week, timeslot, room_online, lecturer, group), 0)
        for week in weeks
        for lecture in lectures
        for lecturer in lecturers
        for group in student_groups
        for timeslot in timeslots
    )
    model.Add(penalty_online_lecture == penalty_online_lecture_usage)
    objective_function_with_penalty_online_lecture = (
        penalty_online_lecture * penalty_weight
    )
    return objective_function_with_penalty_online_lecture

# printing function
def print_schedule(lectures: dict, weeks: list, timeslots: list, rooms: list, lecturers: dict, student_groups: dict, solver: cp_model.CpSolver) -> None:
    """
    Print the Schedule
    """
    for week in weeks:
        print(f"\n{week}")
        for i in range(6):
            match i:
                case 0:
                    print("Monday")
                case 1:
                    print("Tuesday")
                case 2:
                    print("Wednesday")
                case 3:
                    print("Thursday")
                case 4:
                    print("Friday")
                case 5:
                    print("Saturday")
            print(
                tabulate(
                    [
                        [timeslot % 3, room.name, group, lecture, lecturer]
                        for timeslot in timeslots
                        for room in rooms
                        for lecturer in lecturers
                        for group in student_groups
                        for lecture in lectures
                        if lecture in lecturers[lecturer]
                        and lecture in student_groups[group]
                        and solver.Value(
                            schedule[
                                (lecture, week, timeslot, room, lecturer, group)
                            ]
                        )
                        == 1
                        and timeslot < (3 * (i+1))
                        and timeslot >= 3 * i
                    ],
                    headers=["timeslot", "room", "group", "lecture", "lecturer"],
                    tablefmt="orgtbl",
                    floatfmt=".1f"
                )
            )

# Create the model.
model = cp_model.CpModel()

weeks = create_weeks(12)
timeslots = create_timeslots(18)

# Weights
penalty_weight_first_block = 10
penalty_weight_third_block = 15
penalty_weight_saturday_block = 15
penalty_weight_high_ranking = 5
penalty_weight_last_three_weeks = 20
penalty_weight_online_lecture = 5

#RAUM ONLINE MUSS AM LETZTEN INDEX DES ARRAYS SEIN!!!! SONST KAPUTTI
rooms = []
raum_4_0 = Room("4.0", 4, 20, 1)
raum_4_1 = Room('4.1', 4, 36, 2)
raum_4_2 = Room('4.2', 4, 36, 3)
raum_4_3 = Room('4.3', 4, 36, 4)
raum_online = Room("online", 0, 99999, 99999)

rooms.append(raum_4_0)
rooms.append(raum_4_1)
rooms.append(raum_4_2)
rooms.append(raum_4_3)
rooms.append(raum_online)

json_struct = []

with open("dits1.json", "r") as file:
    json_struct.append(json.load(file))

with open("dits2.json", "r") as file:
    json_struct.append(json.load(file))

with open("dwis1.json", "r") as file:
    json_struct.append(json.load(file))
    
with open("dwis2.json", "r") as file:
    json_struct.append(json.load(file))

student_groups = {}
vetos = {}
lecturers = {}
lectures = {}

for i in range(len(json_struct)):
    student_groups.update(json_parser.getGroup(json_struct[i]))
    lecturers.update(json_parser.getLecturers(json_struct[i]))
    lectures.update(json_parser.getLectures(json_struct[i]))

vetos = {
    " Lobachev": {
        "KW1": [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17]
    },
    " Werner": {"KW5": [6, 7, 8, 9, 10]},
}

schedule = create_schedule_vars(lectures, weeks, timeslots, rooms, lecturers, student_groups, model)

## add hard constraints
add_each_lecture_is_assigned_hard_constraint(lectures, weeks, timeslots, rooms, lecturers, student_groups, model)
add_lectures_can_not_be_in_same_timeslot_and_room_in_the_same_week_hard_constraint(lectures, weeks, timeslots, rooms, lecturers, student_groups, model)
add_lecturer_can_not_be_assigned_twice_hard_constraint(lectures, weeks, timeslots, rooms, lecturers, student_groups, model)
add_student_group_can_not_be_assigned_twice_hard_constraint(lectures, weeks, timeslots, rooms, lecturers, student_groups, model)
add_lecturers_can_veto_days_hard_constraint(lectures, weeks, timeslots, rooms, lecturers, student_groups, model)
add_online_and_on_site_lecture_must_not_be_on_same_day_per_group_hard_constraint(lectures, weeks, timeslots, rooms, lecturers, student_groups, model)
add_max_30_percent_online_per_group_hard_constraint(lectures, weeks, timeslots, rooms, lecturers, student_groups, model)

## add soft constraints
objective_function_with_penalty_first = add_minimize_first_block_per_day_soft_constraint(penalty_weight_first_block, lectures, weeks, timeslots, rooms, lecturers, student_groups, model)
objective_function_with_penalty_third = add_minimize_third_block_per_day_soft_constraint(penalty_weight_third_block, lectures, weeks, timeslots, rooms, lecturers, student_groups, model)
objective_function_with_penalty_saturday = add_minimize_saturday_lectures_soft_constraint(penalty_weight_saturday_block, lectures, weeks, timeslots, rooms, lecturers, student_groups, model)
objective_function_with_penalty_high_ranking = add_minimize_rooms_with_high_ranking_soft_constraint(penalty_weight_high_ranking, lectures, weeks, timeslots, rooms, lecturers, student_groups, model)
objective_function_with_penalty_last_three_weeks = add_minimize_lectures_in_last_three_weeks_soft_constraint(penalty_weight_last_three_weeks, 3, lectures, weeks, timeslots, rooms, lecturers, student_groups, model)
objective_function_with_penalty_online_lecture = add_minimize_online_lectures_soft_constraint(penalty_weight_online_lecture, lectures, weeks, timeslots, raum_online, lecturers, student_groups, model)


model.Minimize(
    objective_function_with_penalty_third
    + objective_function_with_penalty_first
    + objective_function_with_penalty_saturday
    + objective_function_with_penalty_last_three_weeks
    + objective_function_with_penalty_high_ranking
    + objective_function_with_penalty_online_lecture
)  

# Solve the model.
# We use Google OR-Tools to solve the model.
solver = cp_model.CpSolver()
status = solver.Solve(model)



# Print the schedule.
# If the model is optimal, we print the optimal schedule.
if status == cp_model.OPTIMAL:
    print_schedule(lectures, weeks, timeslots, rooms, lecturers, student_groups, solver)
else:
    print(f"Das ist mit den Parametern nicht lösbar, du Minion!")