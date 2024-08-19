<script lang="ts">
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

		SvelteKit Todo app
	</h1>

	<table>
		<thead>
			<tr>
				<th class="taskId">id</th>
				<th>Task</th>
				<th>Status</th>
				<th>Delete</th>
			</tr>
		</thead>
		<tbody>
			{#each tasks as task}
				<tr>
					<td>{task.id}</td>
					<td>{task.name}</td>
					<td class="task-status"
						><button
							class="task-status-button"
							type="button"
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
							}}
						>
							{task.status == Status.DONE ? "✔" : "❌"}
						</button></td
					>
					<td class="task-del">
						<button
							class="task-status-button"
							type="button"
							on:click={async () => {
								if (
									!confirm(
										`Are you sure to delete task ${task.id} ${task.name}?`,
									)
								) {
									return;
								}

								await client.deleteTask({
									id: task.id,
								});

								systemMessage = `Task ${task.id} is deleted.`;
								getAllTasks();
							}}
						>
							🗑️
						</button>
					</td>
				</tr>
			{/each}
		</tbody>
	</table>

	<form
		class="task-form"
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
		<div class="input-area">
			<p class="form-label">タスク名</p>
			<input class="form-text-input" bind:value={inputValue} />
			<button class="form-submit-button" type="submit">追加</button>
		</div>
	</form>
	<p class="system-message">{systemMessage}</p>
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

	.task-status-button {
		min-width: 4em;
		min-height: 2em;
		border: 0;
		background: transparent;
		border-radius: 5px;
		cursor: pointer;
	}

	.task-form {
		box-sizing: border-box;
		margin-top: 0.5em;
		padding: 0 1em;
		width: 100%;
	}

	.input-area {
		width: 100%;
		display: flex;
		justify-content: space-between;
		justify-items: center;
		align-items: center;
		gap: 0.5em;
	}

	.form-label {
		padding: 0;
		margin: 0;
		flex-basis: 6em;
	}

	.form-text-input {
		flex-basis: 100%;
		padding: 0.5em;
	}

	.form-submit-button {
		flex-basis: 5em;
		appearance: none;
		border: 0;
		border-radius: 5px;
		background: #4676d7;
		color: #fff;
		padding: 8px 16px;
		font-size: 16px;
		cursor: pointer;
	}

	.form-submit-button:hover {
		background: #1d49aa;
	}

	.form-submit-button:focus {
		outline: none;
		box-shadow: 0 0 0 4px #cbd6ee;
	}
</style>
