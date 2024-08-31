import os
import sys

source = ''
def insertTxt():
    return '''
import "gogoproto/gogo.proto";

option go_package                           = "%s";
option (gogoproto.goproto_enum_prefix_all)  = false;
option (gogoproto.goproto_unrecognized_all) = false;
option (gogoproto.goproto_unkeyed_all)      = false;
option (gogoproto.goproto_sizecache_all)    = false;
''' % (source)

def gogoFile():
    if os.path.exists(source + '/tmp'):
        os.system('rm -r {}/tmp'.format(source))
    proto_files = []
    for root, dirs, files in os.walk(source):
        if root.find('gogoproto') != -1:
            continue
        for file in files:
            if file.endswith('.proto'):
                proto_files.append(os.path.join(root, file))
    os.mkdir(source + '/tmp')
    os.system("cp -r {}/gogoproto {}/tmp".format(sys.path[0], source))
    for file in proto_files:
        with open(file, 'r+') as f:
            txt = f.read()
            index = insertIndex(txt)
            if index == -1:
                continue
            if txt.find('gogoproto') == -1:
                txt = txt[:index] + insertTxt() + txt[index:]
        with open(source + '/tmp/'+file.split('/')[-1], 'w') as f:
            f.write(txt)
    os.system('protoc -I={}/tmp --gofast_out=paths=source_relative:./{}/tmp {}/tmp/*.proto'.format(source, source, source, source))
    os.system('mv {}/tmp/*.go {}/'.format(source, source))
    os.system('rm -r {}/tmp'.format(source))

def insertIndex(txt):
    i1 = txt.find('message')
    i2 = txt.find('enum')
    index = 1E10
    if i1 != -1:
        index = i1
    if i2 != -1:
        index = min(index, i2)
    if index == 1E10:
        return -1
    return index

if __name__ == '__main__':
    if len(sys.argv) < 2:
        print("cmd argv error")
        exit()
    for arg in sys.argv[1:]:
        key, value = arg.split("=")
        if key == "source":
            source = value
    if source == '':
        print("path is empty")
        exit()
    gogoFile()
    print("generated successfully".format(source))
