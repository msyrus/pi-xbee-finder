import React, { Component } from 'react';

import { Rect, Circle, Wedge, Line, Stage, Layer, Group } from "react-konva";

export class OfficeMap extends Component {
	state = {
		color: "green",
	};

	constructor(props) {
		super(props)
		this.baseUrl = props.baseUrl;
		this.state = {
			dev1: {
				key: props.dev1,
				value: false,
			},
			dev2: {
				key: props.dev2,
				value: false,
			},
			dev3: {
				key: props.dev3,
				value: false,
			},
		}
	}

	componentDidMount() {
		this.fetchTimer = setInterval(
			() => this.fetchData(),
			1000
		);
	}

	componentWillUnmount() {
		clearInterval(this.fetchTimer)
	}

	fetchData() {
		fetch(this.baseUrl)
		.then(res => res.json())
		.then(
			(res) => {
				let state = {};
				for (let dev in res.data ) {
					if ( dev === this.state.dev1.key) {
						state["dev1"] = {key: dev, value: res.data[dev]};
					}
					if ( dev === this.state.dev2.key) {
						state["dev2"] = {key: dev, value: res.data[dev]};
					}
					if ( dev === this.state.dev3.key) {
						state["dev3"] = {key: dev, value: res.data[dev]};
					}
				}
				if (res.error) {
					console.log(res.error);
				}
				this.setState(state);
			},
			(error) => {
				console.log(error);
			}
		)
	}

	calcZone(d1, d2, d3) {
		let zone = d1;
		zone += (d2<<1);
		zone += (d3<<2);
		return zone;
	}

	renderZone0() {
		return (
			<Group>
			</Group>
		);
	}

	renderZone1() {
		return (
			<Group>
				<Wedge
					fill={"gray"}
					x={170}
					y={165}
					radius={300}
					angle={180}
				/>
				<Wedge
					fill={"white"}
					x={660}
					y={165}
					radius={300}
					angle={180}
					globalCompositeOperation={"destination-out"}
				/>
				<Wedge
					fill={"white"}
					x={430}
					y={20}
					radius={250}
					angle={50}
					rotation={63.5}
					globalCompositeOperation={"destination-out"}
				/>
			</Group>
		);
	}

	renderZone2() {
		return (
			<Group>
				<Wedge
					fill={"gray"}
					x={660}
					y={165}
					radius={300}
					angle={180}
				/>
				<Wedge
					fill={"white"}
					x={170}
					y={165}
					radius={300}
					angle={180}
					globalCompositeOperation={"destination-out"}
				/>
				<Wedge
					fill={"white"}
					x={430}
					y={20}
					radius={250}
					angle={50}
					rotation={63.5}
					globalCompositeOperation={"destination-out"}
				/>
			</Group>
		);
	}

	renderZone3() {
		return (
			<Group>
				<Wedge
					fill={"gray"}
					x={170}
					y={165}
					radius={300}
					angle={180}
					/>
				<Wedge
					fill={"gray"}
					x={660}
					y={165}
					radius={300}
					angle={180}
					globalCompositeOperation={"destination-in"}
				/>
				<Wedge
					fill={"white"}
					x={430}
					y={20}
					radius={250}
					angle={50}
					rotation={63.5}
					globalCompositeOperation={"destination-out"}
				/>
			</Group>
		);
	}

	renderZone4() {
		return (
			<Group clipX={20} clipY={20} clipWidth={800} clipHeight={145}>
				<Wedge
					fill={"gray"}
					x={430}
					y={20}
					radius={250}
					angle={180}
				/>
			</Group>
		);
	}

	renderZone5() {
		return (
			<Group>
				<Wedge
					fill={"gray"}
					x={170}
					y={165}
					radius={300}
					angle={180}
				/>
				<Wedge
					fill={"gray"}
					x={430}
					y={20}
					radius={250}
					angle={50}
					rotation={63.5}
					globalCompositeOperation={"destination-in"}
				/>
				<Wedge
					fill={"white"}
					x={660}
					y={165}
					radius={300}
					angle={180}
					globalCompositeOperation={"destination-out"}
				/>
			</Group>
		);
	}

	renderZone6() {
		return (
			<Group>
				<Wedge
					fill={"gray"}
					x={660}
					y={165}
					radius={300}
					angle={180}
					/>
				<Wedge
					fill={"gray"}
					x={430}
					y={20}
					radius={250}
					angle={50}
					rotation={63.5}
					globalCompositeOperation={"destination-in"}
					/>
				<Wedge
					fill={"white"}
					x={170}
					y={165}
					radius={300}
					angle={180}
					globalCompositeOperation={"destination-out"}
				/>
			</Group>
		);
	}

	renderZone7() {
		return (
			<Group>
				<Wedge
					fill={"gray"}
					x={660}
					y={165}
					radius={300}
					angle={180}
					/>
				<Wedge
					fill={"gray"}
					x={430}
					y={20}
					radius={250}
					angle={50}
					rotation={63.5}
					globalCompositeOperation={"destination-in"}
					/>
				<Wedge
					fill={"gray"}
					x={170}
					y={165}
					radius={300}
					angle={180}
					globalCompositeOperation={"destination-in"}
				/>
			</Group>
		);
	}

	renderZone(zone) {
		switch (zone) {
			case 1:
				return this.renderZone1();
			case 2:
				return this.renderZone2();
			case 3:
				return this.renderZone3();
			case 4:
				return this.renderZone4();
			case 5:
				return this.renderZone5();
			case 6:
				return this.renderZone6();
			case 7:
				return this.renderZone7();
			default:
				return this.renderZone0();
		}
	}

	render() {
		let {state} = this;
		let zone = this.calcZone(state.dev1.value, state.dev2.value, state.dev3.value);
		
		return (
			<Stage width={window.innerWidth} height={window.innerHeight}>
				<Layer
					clipX={20}
					clipY={20}
					clipWidth={800}
					clipHeight={500}
				>
					{this.renderZone(zone)}	
					<Circle
						fill={state.dev1.value ? "green" : "red"}
						x={170}
						y={170}
						radius={5}
						fillPriority={100}
					/>
					<Circle
						fill={state.dev2.value? "green" : "red"}
						x={660}
						y={170}
						radius={5}
					/>
					<Circle
						fill={state.dev3.value ? "green" : "red"}
						x={430}
						y={25}
						radius={5}
					/>
					<Rect
						x={20}
						y={20}
						width={800}
						height={500}
						stroke={"black"}
						strokeWidth={2}
					/>
					<Line
						points={[20, 160, 370, 160]}
						stroke={"black"}
						strokeWidth={10}
					/>
					<Line
						points={[500, 160, 820, 160]}
						stroke={"black"}
						strokeWidth={10}
					/>
					<Line
						points={[170, 280, 170, 520]}
						stroke={"black"}
						strokeWidth={10}
					/>
				</Layer>
			</Stage>
		);
	}
}