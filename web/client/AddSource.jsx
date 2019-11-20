import React from 'react';

export default class AddSource extends React.Component {
  constructor(props) {
    super(props)

    this.state = {
      input: false,
      value: ''
    }
    this.onClick = this.onClick.bind(this);
    this.onChange = this.onChange.bind(this);
    this.submitSource = this.submitSource.bind(this);
  }

  onClick(event) {
    this.setState({ input: true });
  }

  onChange(event) {
    event.preventDefault();
    this.setState({ value: event.target.value });
  }

  submitSource(event) {
    event.preventDefault();
    this.props.onSubmit(this.state.value);
    this.setState({ 
      input: false,
      value: ""
    });
  }

  render() {
    let form;
    if (this.state.input === true) {
      form = <form type="submit" onSubmit={this.submitSource} >
        <input type="text" value={this.state.value} onChange={this.onChange} className="inputSource"/>
      </form>
    }
    return (
      <div className="addSource">
        <button onClick={this.onClick} className="addButton"> + Add Source</button>
        {form}
      </div>
    )
  }
}