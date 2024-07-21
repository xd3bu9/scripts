# A python script that can extract metadata of a pdf and text from all the pages of a pdf
# Pre-requisite: pip install PyPDF2
# usage: python getpdfinfo.py /path/to/file.pdf

from PyPDF2 import PdfReader
from argparse import ArgumentParser

def get_args():
	parser = ArgumentParser()
	parser.add_argument('FILEPATH', help='/pdf/file/path')
	return parser.parse_args()

def get_info(path):
	with open(path, 'rb') as f:
		pdfFile = PdfReader(f)
		pdfMetadata = pdfFile.metadata
		number_of_pages = len(pdfFile.pages)
		print("--------METADATA--------")
		print(pdfMetadata)
		print("------------------------")
		get_text(pdfFile)

def get_text(file):
	print("--------TEXT--------")
	for page in file.pages:
		print(page.extract_text())
	print("--------------------")

if __name__ == '__main__':
	args = get_args()
	path = args.FILEPATH
	get_info(path)
