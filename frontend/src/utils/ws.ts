import {LogMessage} from "./models";

export function onClose(){
    console.log("closed")
}

export function onError(e: any){
    console.log("error! " + e)
}

export function onOpen(){

}
