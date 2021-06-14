#!/usr/bin/env python3
import csv
import requests
import time


def main():
    print('Starting')
    for publisher, publication, issue in csv.reader(open("items.csv")):
        t0 = time.time()
        requests.post(
            "http://localhost:5000/",
            json={
                "Publisher": int(publisher),
                "Publication": int(publication),
                "Issue": int(issue),
            },
        )
        delta = time.time() - t0
        print('Parsed {}/{}/{} in {} seconds'.format(publisher, publication, issue, delta))


if __name__ == "__main__":
    main()
