#!/usr/bin/env python3

import argparse
import os
import random
import sys
import threading
from itertools import count
from string import ascii_letters, digits

from colorama import Fore
from kahoot import client


def getargs():
    parser = argparse.ArgumentParser(description='Kahoot Flooder Bot - Created by Music_Dude#0001')
    parser.add_argument('code', metavar='CODE', type=int)

    parser.add_argument('-t', '--threads', type=int, default=1, help='number of threads to run in each process')

    group = parser.add_mutually_exclusive_group()
    group.add_argument('-n', '--nick', type=str, default='', metavar='NICK', help='nickname to use, e.g. nick1, nick2...')
    group.add_argument('-r', '--random-nick', dest='rand', action='store_true', help='join using random nicknames')

    return parser.parse_args()

def main():
    global i
    while True:
        if args.rand:
            name = ''.join(random.choices(ascii_letters + digits, k=8))
        else:
            name = args.nick + str(i+1)

        if client().join(args.code, name) is None:
            print(f'{Fore.LIGHTBLUE_EX}{threading.current_thread().getName():10}{Fore.RESET} | {Fore.MAGENTA}{i+1:^3}{Fore.RESET} | Joining as {Fore.RED}{name}{Fore.RESET}')
        else:
            print(f'{Fore.RED}Failed to join {args.code}!{Fore.RESET}')
            exit(1)

        i += 1

if __name__ == '__main__':
    args = getargs()
    threads = []
    i = 0

    for _ in range(args.threads):
        t = threading.Thread(target=main)
        threads.append(t)
    
    for thread in threads:
        thread.start()

