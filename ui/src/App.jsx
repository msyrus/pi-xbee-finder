import React, { Component } from 'react';
// import './App.css';
import { OfficeMap } from './officemap';

class App extends Component {
	state = {
		color: "green",
	};

	render() {
		return (
			<OfficeMap
				dev1={"0013a20040a9c7d4"}
				dev2={"0013a20040a9c870"}
				dev3={"0013a20040a9c7e5"}
				baseUrl={"http://localhost:5678/"}
			/>
		);
	}
}

export default App;
