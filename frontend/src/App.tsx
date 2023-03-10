import React, {useState} from 'react';
import './App.css';
import 'bootstrap/dist/css/bootstrap.min.css';
import {Button, Form} from "react-bootstrap";


function onMessage(this: WebSocket, ev: MessageEvent<any>): void {
    console.log(ev)
}
function onClose(){
    console.log("closed")
}

function onError(e: any){
    console.log("error! " + e)
}

function onOpen(){

}
function connect(text: string): void {
    let ws = new WebSocket(`ws://localhost:8080/ws?message=${encodeURIComponent(text)}`);
    ws.onopen = onOpen;
    ws.onmessage = onMessage;
    ws.onerror = onError;
    ws.onclose = onClose;
}
function App() {
    const [text, setText] = useState("");

    return (
        <Form>
            <Form.Control type="text" placeholder="Enter text" onChange={event => setText(event.target.value)}/>
            <Button variant="primary" type="button" onClick={() => connect(text)}>
                Submit
            </Button>
        </Form>
    );
}

export default App;
