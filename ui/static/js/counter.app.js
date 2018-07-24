class Counter extends React.Component{
    constructor(props){
        super(props)
        this.state = {
            ticketNumber: "0000",
            queueId: "",
            counterId: ""
        }
    }

    render(){
        return(
            <div className="counter-container">
                <div className="ticket-display">
                    <div>
                        <h1 className="display-1 text-center">TICKET NO.  {this.state.ticketNumber}</h1>
                        <h1 className="display-4 text-center">PROCEED TO</h1>
                    </div>
                </div>
                <div className="counter-display">
                    <div>
                        <h1 className="display-1">COUNTER {this.state.counterId}</h1>
                    </div>
                    
                </div>
            </div>
        );
    }

    componentDidMount(){
        var queryParams = new URLSearchParams(location.search);
        var queueId = queryParams.get("queueId");
        var counterId = queryParams.get("counterId");

        this.setState({queueId: queueId, counterId: counterId});

        var ws = new WebSocket("ws://localhost:5000/ws");

        ws.onopen = () => {
            ws.send(JSON.stringify({
                "counterId": this.state.counterId,
                "queueId": this.state.queueId
            }));
        }
        
        ws.onmessage = (msg) =>{
            var payload = JSON.parse(msg.data)
            if (payload.type === "update"){
                this.setState({ticketNumber: payload.data})
            }
        };
    }
}

ReactDOM.render(<Counter />, document.getElementById("root"));