import time, typing, uuid
import pydantic


class PutReport(pydantic.BaseModel):
    reporter: str
    reason: str

    reported: str = None
    id: str = str(uuid.uuid4())
    created: int = int(time.time())
    resolved: bool = False

    class Config:
        orm_mode = True


class Report(pydantic.BaseModel):
    reporter: str
    reason: str

    reported: str
    id: str
    created: int
    resolved: bool
