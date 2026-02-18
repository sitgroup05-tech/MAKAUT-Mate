

import os
import argparse
import dotenv

import types

from enum import Enum
from collections import namedtuple

import psutil
import pipe
from typing_extensions import Self , NamedTuple

dotenv.load_dotenv()

# RPC
def AddRPCTasks(RPCTaskParser : argparse.ArgumentParser) :
    # health-check
    RPCTaskParser.add_argument("--is-alive",action="store_true",default=False,help="is the rpc server is running or not")
    # details 
    RPCTaskParser.add_argument("--details",action="store_true",help="logs out details of the server.")
 
    RPCTasksubParser = RPCTaskParser.add_subparsers(dest="rpc_subcommand")

    # Manual-server start
    RPCServerTask = RPCTasksubParser.add_parser("start-server",help="start the server")
    RPCServerTask.add_argument("-H","--host",help="Host",default=os.getenv("RPC_HOST") or "localhost")
    RPCServerTask.add_argument("-P","--port",help="Port",type=int,default=int(os.getenv("RPC_PORT") or "8080"))


# Admin
def AddAdminTasks(AdminParser : argparse.ArgumentParser):

    AdminParser.add_argument("--list",help="list out the admin user details",action="store_true")

    AdminUserParser = AdminParser.add_subparsers(dest="admin_subcommand")

    # add
    AdminUserAppend = AdminUserParser.add_parser("append",help="appends a new admin user")
    AdminUserAppend.add_argument("name",help="name")
    AdminUserAppend.add_argument("mail",help="E-mail")
    AdminUserAppend.add_argument("password",help="password")

    # remove 
    AdminUserRemove = AdminUserParser.add_parser("remove",help="removes user from network")



class Commands(Enum) :
    RPC = "rpc"
    ADMIN = "admin"

class Result(types.SimpleNamespace) :
    Arg : NamedTuple


class RPCManager :

    def __init__(self : Self , args : argparse.Namespace ) : 
        
        should_start_server = args.rpc_subcommand == "start-server"
                
        self.task = (
            Result() 
                | pipe.Pipe(lambda x : self._start_server(x) if should_start_server else x ) 
                | self._check_is_alive
                | pipe.Pipe(lambda x : self._log_details(x) if args.details else x) 
        )  
        

    @pipe.Pipe
    def _start_server(self : Self , r : Result) -> Result :
        ...

    @pipe.Pipe
    def _log_details(self : Self , r : Result) -> Result :
        if not r.is_alive :
            return r 
        



    @pipe.Pipe
    def _check_is_alive(self : Self , r : Result) -> Result:
        pid = os.getenv("PYTHON_RPC_PID")
        active = os.getenv("PYTHON_RPC_ACTIVE",default="")

        if active == "True" :
            active =  psutil.pid_exists(int(pid))
            if not active :
                path = dotenv.find_dotenv()
                dotenv.set_key(path,"PYTHON_RPC_ACTIVE","False")
                dotenv.set_key(path,"PYTHON_RPC_PID","")
                r.is_alive = False
            else :  
                r.is_alive = True
        
            return r 
        r.is_alive = False
        return r 

    def __call__(self):
        self.task()
        return self.Result 

class AdminManager :

    def __init__(self):
        
        pass

    def __call__(self, *args, **kwds):
        pass


if __name__ == "__main__":

    parser = argparse.ArgumentParser(prog="",description="",epilog="")
    subParser = parser.add_subparsers(dest="command")

    AddRPCTasks(
        subParser.add_parser(Commands.RPC,help="this offers service related to the the PY-RPC")
    )

    AddAdminTasks(
        subParser.add_parser(Commands.ADMIN,help="admin level user task")
    )

    args = parser.parse_args()

    match args.command : 
        case Commands.RPC : 
            manager = RPCManager(args)
        case Commands.ADMIN:
            manager = AdminManager(args)
        case _ :
            raise AssertionError("Unreachable Code")
    
    manager()

    print(manager)
