<script lang="ts">
    import { liberateState } from "./liberateStore.svelte";
    import { fade } from "svelte/transition";

    let hasResults = $derived(liberateState.links.length > 0);
</script>

{#if hasResults}
    <div
        class="results-container flex flex-col gap-3 max-w-xl w-full px-8 mt-20"
        transition:fade={{ duration: 200 }}
    >
        {#each liberateState.links as link}
            <a
                class="card p-2 preset-tonal-primary hover:preset-tonal-surface transition-colors block"
                href={link}
                target="_blank"
                rel="noopener noreferrer"
            >
                <div class="text-secondary-500 hover:underline break-all">
                    {link}
                </div>
            </a>
        {/each}
    </div>
{/if}
