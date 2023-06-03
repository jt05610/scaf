package py

import "text/template"

func main() *template.Template {
	t, err := template.New("main.py").Parse(`
from fastapi import FastAPI

app = FastAPI()


@app.get("/")
async def root():
    return {"message": "Hello World"}


@app.get("/hello/{name}")
async def say_hello(name: str):
    return {"message": f"Hello {name}"}

`)
	if err != nil {
		panic(err)
	}
	return t
}

func testMain() *template.Template {
	t, err := template.New("test_main.http").Parse(`
GET http://{{.Addr}}:{{.Port}}/
Accept: application/json

###

GET http://{{.Addr}}:{{.Port}}/hello/User
Accept: application/json

###
`)
	if err != nil {
		panic(err)
	}
	return t
}

func Install() *template.Template {
	t, err := template.New("install.sh").Parse(`python -m pip install virtualenv
virtualenv {{.ParDir}}/venv
source {{.ParDir}}/venv/bin/activate
python -m pip install --upgrade pip
python -m pip install grpcio
python -m pip install fastapi
python -m pip install uvicorn
`)
	if err != nil {
		panic(err)
	}
	return t
}

func Serve() *template.Template {
	t, err := template.New("serve.sh").Parse(`source {{.ParDir}}/venv/bin/activate
uvicorn main:app --reload --port {{.Port}} --host {{.Addr}}
`)
	if err != nil {
		panic(err)
	}
	return t
}
