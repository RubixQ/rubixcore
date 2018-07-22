class Modal extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            display: false,
            confirmWasCalled: false,
        };
    }

    componentDidMount() {
        var self = this;

        this.$modalElem = $('#app-modal');

        this.$modalElem.on('hidden.bs.modal', e => {
            if (self.props.type === 'confirm' && !self.state.confirmWasCalled) {
                self.props.confirm(false, 'closeEvent');
            }

            self.props.onClose && self.props.onClose();
            this.setState({ display: false, confirmWasCalled: false });
        });

        this.$modalElem.on('show.bs.modal', e => {
            self.setState({ display: true });
        });

        if (this.props.display) {
            this.$modalElem.modal('show');
        }
    }

    componentDidUpdate() {
        if (this.props.display) {
            return this.$modalElem.modal('show');
        }

        return this.$modalElem.modal('hide');
    }

    render() {
        const { content, title, type } = this.props;
        const self = this;

        return (
            <div className="modal fade in" id="app-modal" tabIndex="-1" role="dialog">
                <div className="modal-dialog" role="document">
                    <div className="modal-content">
                        <div className="modal-header">
                            <h5 className="modal-title" id="exampleModalLabel">{title}</h5>
                            <button type="button" className="close" data-dismiss="modal" aria-label="Close">
                                <span aria-hidden="true">&times;</span>
                            </button>
                        </div>

                        <div className="modal-body">{content}</div>

                        <div className="modal-footer">
                            {type === 'info' ? (
                                <button
                                    type="button"
                                    className="btn btn-default waves-effect"
                                    data-dismiss="modal"
                                >
                                    Close
                                </button>
                            ) : (
                                    <React.Fragment>
                                        <button
                                            className="btn btn-default wavonCloseect"
                                            data-dismiss="modal"
                                            onClick={evt => {
                                                self.props.confirm(false);
                                                self.setState({ confirmWasCalled: true });
                                            }}
                                        >
                                            {this.props.closeButton || 'Close'}
                                        </button>
                                        <button
                                            className="btn btn-info waves-effect waves-light"
                                            onClick={evt => {
                                                self.setState({ confirmWasCalled: true });
                                                self.props.confirm(true);
                                            }}
                                        >
                                            {this.props.okButton || 'OK'}
                                        </button>
                                    </React.Fragment>
                                )}
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}


const FormInput = ({ value, onChange }) => {
    return (
        <div className="form-group">
            <label htmlFor="phoneNumber">Enter phone number to receive ticket number via sms</label>
            <input type="text" value={value} onChange={onChange} className="form-control" id="phoneNumber" placeholder="Enter phone number (02xxxxxxxx)" />
        </div>
    );
}

class Kiosk extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            queues: [],
            phoneNumber: '',
            hasError: false,
            showModal: false,
            clickedQueueId: '',
            queueName: ''
        };
    }

    setQueueDetails = (id, name) => {
        this.setState({ clickedQueueId: id, queueName: name });
    }

    toggleModal = () => {
        this.setState({ showModal: !this.state.showModal });
    }

    handleChange = ({ target }) => {
        this.setState({ phoneNumber: target.value });
    }

    confirmSend = (value) => {
        if (value) {
            axios.post("http://192.168.178.62:5000/customers", { msisdn: this.state.phoneNumber, queueId: this.state.clickedQueueId, queueName: this.state.queueName })
            .then((res) => {
                console.log(res);
            })
            .catch((err) => {
                this.setState({ hasError: true });
            });
        }
    };

    render() {
        return (
            <div className="queues-container">
                <div className="title">
                    <h1>WELCOME</h1>
                    <p>Please select an action</p>
                </div>
                {this.state.queues.length ? this.state.queues.map(q =>
                    <button type="button" key={q.id} className="btn action-items btn-primary btn-lg btn-block" onClick={() => {
                        this.setQueueDetails(q.id, q.name);
                        this.toggleModal();
                    }}>{q.name}</button>
                ) : <p>Loading....</p>}
                <Modal confirm={this.confirmSend.bind(this)} display={this.state.showModal} onClose={this.toggleModal.bind(this)} title='Enter phone number' type='confirm' content={<FormInput onChange={this.handleChange.bind(this)} value={this.state.phoneNumber} />} />
            </div>
        );
    }

    componentDidMount() {
        axios.get("http://192.168.178.62:5000/queues")
            .then((res) => {
                this.setState({ queues: res.data.data })
                console.log(this.state.queues)
            })
            .catch((err) => {
                this.setState({ hasError: true });
            });
    }
}

ReactDOM.render(<Kiosk />, document.getElementById("root"));