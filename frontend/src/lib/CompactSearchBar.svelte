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
    class="fixed top-4 left-4 flex items-center gap-3 z-50"
    transition:fade={{ duration: 200 }}
>
    <div
        class="text-2xl font-bold bg-gradient-to-r from-primary-500 to-primary-600 bg-clip-text text-transparent"
    >
        Liberate
    </div>
    <input
        class="input flex-1 max-w-2xl"
        bind:value={liberateState.prompt}
        type="text"
        placeholder="Ponder..."
        onkeydown={handleKeydown}
        oninput={handleInput}
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
