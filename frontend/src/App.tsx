import React, {useEffect, useRef, useState} from 'react';
import './App.css';
import 'bootstrap/dist/css/bootstrap.min.css';
import {onOpen, onError, onClose} from "./utils/ws";
import {LogMessage, Problem} from "./utils/models";
import {Button, Container, Grid} from "@mui/material";

let ws




function fetchProblems(): Problem[] {
    return [
        new Problem("problem_fof.p", "problem_fof.p"),
        new Problem("problem_sat.p", "problem_sat.p"),
        new Problem("problem.p", "problem.p"),
        new Problem("PUZ001-1.p", "PUZ001-1.p"),
    ]
}

function cancelLaunch(): void {
    console.log("cancelling launch")
}


function ProblemCard({problem, setMessages, isActive}: { problem: Problem, setMessages: Function, isActive: boolean }) {
    let messages: LogMessage[] = []
    function onMessage(this: WebSocket, ev: MessageEvent<any>): void {
        const parsedJson = JSON.parse(ev.data);
        const msg = new LogMessage(parsedJson.time as Date, parsedJson.source as string, parsedJson.text as string)
        messages = messages.slice()
        messages.push(msg)
        setMessages(messages)
    }
    function launch(problemId: string): void {
        console.log("launching problem " + problemId)
        ws = new WebSocket(`ws://localhost:8080/launch?problemId=${encodeURIComponent(problemId)}`);
        ws.onopen = onOpen;
        ws.onmessage = onMessage;
        ws.onerror = onError;
        ws.onclose = onClose;
    }

    return (
        <Grid className="problem-card" justifyContent="space-between" direction="row" container>
            <Grid item>
                <div className="problem-name">Problem {problem.name}</div>
            </Grid>
            <Grid item>
                <Button className="problem-launch" variant="outlined" onClick={() => launch(problem.id)}>
                    launch
                </Button>
            </Grid>
        </Grid>

    )
}
function ProblemsPanel({setMessages}: {setMessages: Function}){
    const problems = fetchProblems()
    return (
        <Grid container className='problems-panel'  spacing={1} direction='column' columns={2}>
            {problems.map(p =>
                <Grid item xs={6} key={p.id}>
                    <ProblemCard problem={p} setMessages={setMessages} isActive={false}/>
                </Grid>
            )}
        </Grid>
    );
}
function TerminalWindow({messages}: { messages: LogMessage[]})  {
    const containerRef = useRef<HTMLDivElement>(null);


    useEffect(() => {
        const container = containerRef.current;
        if (container) {
            container.scrollTop = container.scrollHeight;
        }
    }, [messages]);
    return (
        <div className="terminal-window" ref={containerRef}>
            <div className="terminal-content">
                {messages.map((msg, index) => (
                    <div key={index} className="terminal-record">
                        {msg.text}
                    </div>
                ))}
            </div>
        </div>
    );
};

// interface ScrollableListProps {
//     list: string[];
// }
//
// const ScrollableList: React.FC<ScrollableListProps> = ({ list }) => {
//     const containerRef = useRef<HTMLDivElement>(null);
//
//     useEffect(() => {
//         const container = containerRef.current;
//         if (container) {
//             container.scrollTop = container.scrollHeight;
//         }
//     }, [list]);
//
//     return (
//         <div ref={containerRef} style={{ height: "300px", overflowY: "scroll" }}>
//             {list.map((item, index) => (
//                 <div key={index}>{item}</div>
//             ))}
//         </div>
//     );
// };



function App() {
    const [messages, setMessages] = useState([]);
    return (
        <Grid container className='problems-panel'  spacing={1} direction='column' columns={1}>
                <Grid item xs={12}>
                    <ProblemsPanel setMessages={setMessages}/>
                </Grid>
            <Grid item xs={12}>
                <TerminalWindow messages={messages}/>
            </Grid>

        </Grid>
    )

}

export default App;
