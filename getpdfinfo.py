# script to extract metadata and text from all pages in a pdf document
# pip install PyPDF2 before runnig this script
# usage: python getpdfinfo.py /home/user/Downloads/file.pdf

from PyPDF2 import PdfReader
from argparse import ArgumentParser

def get_args():
	parser = ArgumentParser()
	parser.add_argument('FILEPATH', help='/pdf/file/path')
	return parser.parse_args()

def get_info(path):
	with open(path, 'rb') as f:
		pdfFile = PdfReader(f)
		info = pdfFile.metadata
		number_of_pages = len(pdfFile.pages)
		print("--------METADATA--------")
		print(info)
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