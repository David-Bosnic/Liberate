<script lang="ts">
    import { fade } from "svelte/transition";
    import { CallLiberate } from "./api";
    import { liberateState } from "./liberateStore.svelte";

    let loading = $state(false);
    let error = $state<string | null>(null);

    async function handleKeydown(e: KeyboardEvent) {
        if (e.key === "Enter" && liberateState.prompt.trim() && !loading) {
            loading = true;
            error = null;
            try {
                const result = await CallLiberate(liberateState.prompt);
                liberateState.links = result.Links;
            } catch (err) {
                error =
                    err instanceof Error ? err.message : "Something went wrong";
            } finally {
                loading = false;
            }
        }
    }
</script>

<div
    class=" flex flex-col items-center justify-center gap-4 py-10"
    transition:fade={{ duration: 200 }}
>
    <div class="text-7xl text-primary-500 font-semibold">Liberate</div>
    <input
        class="input w-96"
        bind:value={liberateState.prompt}
        type="text"
        placeholder="Ponder..."
        onkeydown={handleKeydown}
        disabled={loading}
    />

    {#if loading}
        <div
            class="size-6 rounded-full border-2 border-secondary-500/30 border-t-secondary-500 animate-spin"
        ></div>
    {/if}

    {#if error}
        <div class="text-error-500">{error}</div>
    {/if}
</div>
