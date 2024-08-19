<script lang="ts">
	import Counter from "./Counter.svelte";
	import welcome from "$lib/images/svelte-welcome.webp";
	import welcome_fallback from "$lib/images/svelte-welcome.png";

	import { onMount } from "svelte";

	import { createPromiseClient } from "@connectrpc/connect";
	import { createConnectTransport } from "@connectrpc/connect-web";

	import { TodoService } from "../gen/todo/v1/todo_connect";
	import { Status, type TodoItem } from "../gen/todo/v1/todo_pb";

	let tasks: TodoItem[] = [];
	let inputValue = "";
	let systemMessage = "";

	const transport = createConnectTransport({
		baseUrl: "http://localhost:8080",
	});

	const client = createPromiseClient(TodoService, transport);

	onMount(() => {
		console.info("main page onMount");
		getAllTasks();
	});

	async function getAllTasks() {
		const res = await client.getAllTasks({});
		tasks = res.items.concat().sort((a, b) => Number(a.id) - Number(b.id));
	}
</script>

<svelte:head>
	<title>Home</title>
	<meta name="description" content="Svelte demo app" />
</svelte:head>

<section>
	<h1>
		<span class="welcome">
			<picture>
				<source srcset={welcome} type="image/webp" />
				<img src={welcome_fallback} alt="Welcome" />
			</picture>
		</span>

		to your new<br />SvelteKit app
	</h1>

	<h2>
		try editing <strong>src/routes/+page.svelte</strong>
	</h2>

	<table>
		<thead>
			<tr>
				<th class="taskId">id</th>
				<th>Task</th>
				<th>Status</th>
			</tr>
		</thead>
		<tbody>
			{#each tasks as task}
				<tr>
					<td>{task.id}</td>
					<td>{task.name}</td>
					<td
						class="task-status"
						on:click={async () => {
							const newStatus =
								task.status == Status.DONE
									? Status.TODO
									: Status.DONE;
							await client.updateTaskStatus({
								id: task.id,
								status: newStatus,
							});

							systemMessage = `Task ${task.id} is ${task.status == Status.DONE ? "DONE" : "TODO"}.`;
							getAllTasks();
						}}>{task.status == Status.DONE ? "✔" : "❌"}</td
					>
				</tr>
			{/each}
		</tbody>
	</table>
	<form
		on:submit={async (e) => {
			e.preventDefault();
			const addTaskName = inputValue.trim();
			inputValue = "";
			if (addTaskName === "") {
				systemMessage = "Task name is empty";
				return;
			}

			const res = await client.createTask({
				name: addTaskName,
				status: Status.TODO,
			});

			getAllTasks();
		}}
	>
		<input bind:value={inputValue} />
		<button type="submit">Send</button>
	</form>
	<p class="system-message">{systemMessage}</p>
	<Counter />
</section>

<style>
	section {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		flex: 0.6;
	}

	h1 {
		width: 100%;
	}

	.welcome {
		display: block;
		position: relative;
		width: 100%;
		height: 0;
		padding: 0 0 calc(100% * 495 / 2048) 0;
	}

	.welcome img {
		position: absolute;
		width: 100%;
		height: 100%;
		top: 0;
		display: block;
	}

	table {
		width: 100%;
		text-align: center;
		border-collapse: collapse;
		border-spacing: 0;
		border-top: solid 1px #778ca3;
	}
	table tr:nth-child(2n + 1) {
		background: #e9faf9;
	}
	table th,
	table td {
		padding: 10px;
		border-bottom: solid 1px #778ca3;
	}

	table th.taskId {
		width: 6em;
	}
</style>
