import fastapi
import time
from datatypes import Report

app: fastapi.FastAPI = fastapi.FastAPI()


@app.get("/new")
async def new_reports(size: int = 50, offset: int = 0):
    return {
        "reports": [{
            "id": it
        } for it in range(offset, size)],
        "size": size,
    }


@app.get("/id/{id}")
async def read_report(id: str):
    return {"id": id}


@app.patch("/id/{id}")
async def patch_report(id: str, report: Report):
    report.id = id
    return report


@app.get("/user/{id}")
async def read_user_reports(id: str, size: int = 50, offset: int = 0):
    return {
        "reports": [{
            "id": it,
            "reported": id,
        } for it in range(offset, size)],
        "size": size,
    }


@app.put("/user/{id}")
async def put_user_report(id: str, report: Report):
    report.reported = id
    print(report.dict())
    return report
