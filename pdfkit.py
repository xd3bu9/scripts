# A python script that can: 
# extract metadata and text from a pdf file
# merge multiple pdf files
# password protect a pdf file

# required: pip install PyPDF2
# usage: python pdfkit.py /path/to/file.pdf
import json
from PyPDF2 import PdfReader
from argparse import ArgumentParser

def get_args():
	parser = ArgumentParser()
	parser.add_argument('-i', metavar="INPUT", help='/path/to/input/pdf/file')
	parser.add_argument('-m', metavar="MODE", help='modes. dump,merge,protect')
	return parser.parse_args()

def dump(file_path):
	pages = []
	pageCount = 0
	metadata = ""
	with open(file_path, 'rb') as f:
		pdf = PdfReader(f)
		metadata = pdf.metadata
		pageCount = len(pdf.pages)
		for page in pdf.pages:
			page_text = page.extract_text()
			pages.append(page_text)
	result = {"metadata": metadata, "page_text": pages, "page_count": pageCount}
	print(json.dumps(result))

def protect(filename):
    out = PyPDF2.PdfWriter()
    open_file = PdfReader(open(filename, "rb"))   
    for i in range(0, len(open_file.pages)):
        out.add_page(open_file.pages[i])
    newfile = open(filename+"_secured.pdf", "wb")
    pwd=input("Enter password:")
    out.encrypt(pwd, use_128bit=True)
    out.write(newfile)
    newfile.close()

def merge():
    pdfiles = []
    for filename in os.listdir('.'):
        if filename.endswith('.pdf'):
            if filename != 'merged.pdf':
                pdfiles.append(filename)
    pdfiles.sort(key = str.lower)
    pdfMerge = PyPDF2.PdfMerger()
    for filename in pdfiles:
        pdfFile = open(filename, 'rb')
        pdfReader = PyPDF2.PdfReader(pdfFile)
        pdfMerge.append(pdfReader)
    pdfFile.close()
    pdfMerge.write('merged.pdf')

if __name__ == '__main__':
	args = get_args()
	path = args.i
	mode = args.m
	if (mode == "dump" and path != None):
		dump(path)
	elif (mode == "protect" and path != None):
		protect(path)
	elif (mode == "merge"):
		merge()
	else:
		print("Invalid options. Try -h")