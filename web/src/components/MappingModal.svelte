<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	interface Mapping {
		path: string;
		url: string;
		mapper: string;
		usecount: number;
	}

	export let isOpen: boolean = false;
	export let mapping: Mapping;
	export let isEdit: boolean = false;

	const dispatch = createEventDispatcher();

	function closeModal() {
		isOpen = false;
	}

	function save() {
		dispatch('save', mapping);
		closeModal();
	}
</script>

{#if isOpen}
	<div class="modal modal-open">
		<div class="modal-box">
			<h3 class="text-lg font-bold">{isEdit ? 'Edit Mapping' : 'Add Mapping'}</h3>
			<div class="form-control">
				<label class="label">Path</label>
				<input
					class="input input-bordered"
					type="text"
					bind:value={mapping.path}
					disabled={isEdit}
				/>
			</div>
			<div class="form-control">
				<label class="label">URL</label>
				<input class="input input-bordered" type="text" bind:value={mapping.url} />
			</div>
			<div class="form-control">
				<label class="label">Mapper</label>
				<input class="input input-bordered" type="text" bind:value={mapping.mapper} />
			</div>
			<div class="modal-action">
				<button class="btn" on:click={closeModal}> Cancel </button>
				<button class="btn btn-primary" on:click={save}> Save </button>
			</div>
		</div>
	</div>
{/if}
