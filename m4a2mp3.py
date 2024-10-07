import time
import os
import glob
from pydub import AudioSegment as convert
from argparse import ArgumentParser


def get_args():
    parser = ArgumentParser()
    parser.add_argument("-i", metavar="INPUT", help="/path/to/file.m4a", default=None)
    parser.add_argument("-b", metavar="BULK", help="bulk convert", default=None, action="store_true")
    return parser.parse_args()


def convert_file(file):
    file_name = file[:-4]
    print("Converting",file_name)
    destination = file_name+".mp3"
    file = convert.from_file(file, format="m4a")
    file.export(destination, format="mp3")
    print("[+] " + destination)


def bulk_convert():
    # find all files that end with m4a in cwd
    files = glob.glob("*.m4a")

    for item in files:
        convert_file(file=item)


def check_path(path):
    if os.path.isdir(path):
        print("invalid m4a file path")
        exit(1)
    elif os.path.isfile(path) and path.endswith("m4a"):
        try:
            convert_file(file=path)
        except Exception as e:
            raise e


args = get_args()

if args.b:
    bulk_convert()
else:
    check_path(args.i)


# print("Single Conversion: 	python m4a2mp3.py -i file.m4a")
# print("Bulk Conversion: 	cd /path/to/folder/with/m4a/files/; python m4a2mp3.py -b")
