import os
import sys

path = ""
name = ""
pkg = ""

def readFile(name):
    with open(name, "r") as f:
        return f.read()


def writeFile(name, text):
    if os.path.exists(name):
        os.remove(name)
    with open(name, "w") as f:
        f.write(text)
    os.system("go fmt " + name)


def genUseCase():
    global path, name
    text = '''
    package %s

    import (
        "eastv2/game/modules/play/domain"
        "eastv2/game/modules/play/domain/api"

        "github.com/tnnmigga/corev2/module/domainops"
    )

    var uc *useCase

    type useCase struct {
        *domain.Domain
    }

    func Init(d *domain.Domain) {
        uc = &useCase{
            Domain: d,
        }
        domainops.RegisterCase[api.I%s](d, domain.%sIndex, uc)
    }

    ''' % (name.lower(), name, name)
    dirname = path + "/domain/impl/" + name.lower()
    writeFile(dirname + "/usecase.go", text)

def genApi():
    global path, name
    text = '''
    package api

    type I%s interface {
    }

    ''' % (name)
    writeFile(path + "/domain/api/" + name.lower() + ".go", text)

def genDomain():
    global path, name
    text = readFile(path + "/domain/domain.go")
    index = text.rfind("caseMaxIndex")
    text = text[:index] + name + "Index\n" + text[index:]
    text += '''
    func (d *Domain) %sCase() api.I%s {
        return d.GetCase(%sIndex).(api.I%s)
    }
    ''' % (name, name, name, name)
    writeFile(path + "/domain/domain.go", text)

def genImpl():
    global path, name
    text = readFile(path + "/domain/impl/impl.go").strip()
    text = text[:-1] + "{}.Init(d)\n".format(name.lower()) + text[-1:]
    index = text.find(")")
    text = text[:index] + "\"{}/{}/domain/impl/{}\"\n".format(pkg, path, name.lower()) + text[index:]
    writeFile(path + "/domain/impl/impl.go", text)

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("cmd argv error")
        exit()
    for arg in sys.argv[1:]:
        key, value = arg.split("=")
        if key == "name":
            name = value
        if key == "path":
            path = value
            if path[-1] == "/":
                path = path[:-1]
    dirname = path + "/domain/impl/" + name.lower()
    mod = readFile("./go.mod")
    mod = mod.strip()
    pkg = mod.split("\n")[0].split(" ")[1].strip()
    if os.path.exists(dirname):
        print("useCase already exists")
        exit()
    os.mkdir(dirname)
    genUseCase()
    genDomain()
    genApi()
    genImpl()
    print(name, "case generated")