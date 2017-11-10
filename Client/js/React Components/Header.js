import React from 'react';
import jquery from 'jquery';

export class Header extends React.Component {

  render() {
    return (
      <div id="header-container"className="container-fluid row">
        <div className="col-lg-1 col-md-1 col-sm-0 col-xs-0" />
        <div id="header" className="col-lg-10 col-md-10 col-sm-12 col-xs-12">
          <div id="name">
            MOES:GBA
          </div>
        </div>
      </div>
    );
  }
}
