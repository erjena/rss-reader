import React from 'react';
import '../public/main.css';

export default class Sources extends React.Component {
  constructor(props) {
    super(props)

    this.onClick = this.onClick.bind(this);
  }

  onClick(event) {
    this.props.onSourceChange(event.target.id);
  }

  render() {
    const sources = this.props.sources.map((source, index) => 
      <div id={index} key={index} className={ source.isChosen ? "sourceItemChosen" : "sourceItem" } onClick={this.onClick}>
        {source.name}
      </div>
    )
    return (
      <div style={{ marginTop: "30px" }}>
        {sources}
      </div>
    )
  }
}
