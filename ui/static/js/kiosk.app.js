class Kiosk extends React.Component{
    
    constructor(props){
        super(props);
        this.state = {
            queues: [],
        };
    }

    render(){
        return (
            <div>
                {this.state.queues.map(q => 
                    <button type="button" key={q.id} className="btn btn-primary btn-lg btn-block">{q.name}</button>
                )}
            </div>
        );
    }

    componentDidMount(){
        axios.get("http://localhost:5000/queues")
        .then((res)=>{
            this.setState({queues: res.data.data})
            console.log(this.state.queues)
        })
        .catch((err)=> {

        });
    }
}

ReactDOM.render(<Kiosk/>, document.getElementById("root"));