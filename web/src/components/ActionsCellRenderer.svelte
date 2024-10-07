<script lang="ts">
	export let data: any;
	export let params: any;

	function onEdit() {
		params.context.openEditModal(data);
	}

	function onDelete() {
		if (confirm('Are you sure you want to delete this mapping?')) {
			params.context.deleteMapping(data.path);
		}
	}
</script>

<div class="flex space-x-2">
	<button class="btn btn-sm btn-primary" on:click={onEdit}> Edit </button>
	<button class="btn btn-sm btn-error" on:click={onDelete}> Delete </button>
</div>
