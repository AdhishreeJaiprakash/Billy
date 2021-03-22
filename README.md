**Install Go**: Instructions [here](https://golang.org/doc/install)

**Run Billy**:
```
cd src && go run .
```

Or 

To generate a binary of this tool, run the command below from src/:
```
go build -o <path-to-file>/<filename>
```

**Instructions**:
1. Start with entering the total bill amount.
2. Then enter the names of the people participating in the bill in the format: ```<name1>,<name2>,<name3>...```
3. Proceed to enter each item in the format: ```<iterm-name>:<price>:<person1>,<person2>,...```
4. If you'd like to delete any entry, enter 'remove'/'r'.
5. If you'd like to list all people, enter 'list people'/'lp'.
6. If you'd like to list all the entries Billy has recorded, enter 'list entry'/'le'.
7. When you'd like Billy to work its magic, enter 'done'/'d'. Follow instructions after.
8. To view instructions again, enter 'print instructions'/'p'.
9. To quit, enter 'quit'/'q'.

**TO-DO**:
- Makefile!
- Provide helper flag for binary file
- Monitor SIGINT
