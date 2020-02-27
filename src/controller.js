import React from 'react';

class RestController extends React.Component {

	constructor(props) {
		super(props);
		this.state = {users: []};
		this.headers = [
			{ key: "user" },
			{ key: "date" },
			{ key: "am_weight" }
		];
	}

	componentDidMount() {
		fetch('http://127.0.0.1:8000/days/')
			.then(response => {
				return response.json();
			}).then(result => {
				this.setState({
					users:result
				});
			});
	}
	render() {
		return (
			<table>
				<thead>
					<tr>
					{
						this.headers.map(function(h) {
							return (
								<th key = {h.key}>{h.key}</th>
							)
						})
					}
					</tr>
				</thead>
				<tbody>
					{
						this.state.users.map(function(item, key) {
						return (
								<tr key = {key}>
								  <td>{item.user}</td>
								  <td>{item.date}</td>
								  <td>{item.am_weight}</td>
								</tr>
							)
						})
					}
				</tbody>
			</table>
		)
	}
}

export default RestController;
