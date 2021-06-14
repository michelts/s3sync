#!/usr/bin/env python3
import csv
import requests


def main():
    for publisher, publication, issue in csv.reader(open("items.csv")):
        print('Parsing {}/{}/{}'.format(publisher, publication, issue))
        requests.post(
            "http://localhost:5000/",
            json={
                "Publisher": int(publisher),
                "Publication": int(publication),
                "Issue": int(issue),
            },
        )


if __name__ == "__main__":
    main()
