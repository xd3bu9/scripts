# A python script that can extract metadata from a pdf and the text from all its pages
# required: pip install PyPDF2
# usage: python getpdfinfo.py /path/to/file.pdf

from PyPDF2 import PdfReader
from argparse import ArgumentParser

def get_args():
	parser = ArgumentParser()
	parser.add_argument('FILEPATH', help='/pdf/file/path')
	return parser.parse_args()

def get_info(path):
	with open(path, 'rb') as f:
		pdf = PdfReader(f)
		metadata = pdf.metadata
		pageCount = len(pdf.pages)
		print("--------METADATA--------")
		print(metadata)
		print("------------------------")
		get_text(pdf)
		print("pages: {}", pageCount)

def get_text(file):
	print("--------TEXT--------")
	for page in file.pages:
		print(page.extract_text() + "\n")
	print("--------------------")

if __name__ == '__main__':
	args = get_args()
	path = args.FILEPATH
	get_info(path)
