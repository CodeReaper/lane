import sys
import csv

def escape(str_xml: str):
    str_xml = str_xml.replace("&", "&amp;")
    str_xml = str_xml.replace("<", "&lt;")
    str_xml = str_xml.replace(">", "&gt;")
    str_xml = str_xml.replace("\"", "\\\"")
    str_xml = str_xml.replace("'", "\\'")
    return str_xml

key = int(sys.argv[1]) - 1
value = int(sys.argv[2]) - 1
should_escape = int(sys.argv[3]) == 1

with open(sys.stdin.fileno()) as file:
    reader = csv.reader(file, delimiter=',')
    for row in reader:
        if row[key]:
            v = row[value].replace("\n", "\\n")
            print("|".join([row[key], escape(v) if should_escape else v]))
