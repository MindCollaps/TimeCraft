class Room:
    name = ""
    floor = 0
    capacity = 0
    ranking = 0
    
    def __init__(self, pName: str, pFloor: int, pCapacity: int, pRanking: int):
      self.name = pName
      self.floor = pFloor
      self.capacity = pCapacity
      self.ranking = pRanking
    
    def __str__(self):
        return f"{self.name}---{self.floor}---{self.capacity}---{self.ranking}"
    
    def getName(self) -> str:
        return f"{self.name}"