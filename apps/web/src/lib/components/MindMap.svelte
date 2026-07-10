<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { Application, Assets, Container, Graphics, Sprite, Text, Texture, Color, BlurFilter } from 'pixi.js';
	import mindmapData from '$lib/data/mindmap.json';
	import mainNeuronImage from '$lib/assets/neurons/main_neuron.svg';
	import leafNeuronImage from '$lib/assets/neurons/leaf_neuron.svg';

	type MindMapSection = {
		id: number;
		name: string;
		phrases: string[];
	};

	type LeafNode = {
		graphic: Graphics;
		glow: Graphics;
		label: Text;
		imageSrc?: string;
		image?: Sprite;
		labelWidth: number;
		labelHeight: number;
		homeX: number;
		homeY: number;
		x: number;
		y: number;
		vx: number;
		vy: number;
		phaseX: number;
		phaseY: number;
	};

	type MainNode = {
		graphic: Graphics;
		glow: Graphics;
		label: Text;
		imageSrc?: string;
		image?: Sprite;
		labelWidth: number;
		labelHeight: number;
		homeX: number;
		homeY: number;
		x: number;
		y: number;
		vx: number;
		vy: number;
		phaseX: number;
		phaseY: number;
	};

	type Cluster = {
		id: number;
		color: number;
		container: Container;
		bubble: Graphics;
		links: Graphics;
		main: MainNode;
		leaves: LeafNode[];
	};

	const sections = mindmapData as MindMapSection[];

	const SCENE_SCALE = 1.4;
	const MAIN_NODE_CIRCLE_SIZE = 25 * SCENE_SCALE;
	const MAIN_NODE_HIGHLIGHT_RADIUS = MAIN_NODE_CIRCLE_SIZE * 3;
	const MAIN_NODE_ALPHA = 1;
	const MAIN_NODE_HIGHLIGHT_ALPHA = 0.28;
	const LEAF_NODE_CIRCLE_SIZE = 15 * SCENE_SCALE;
	const LEAF_NODE_HIGHLIGHT_RADIUS = LEAF_NODE_CIRCLE_SIZE * 1.7;
	const LEAF_NODE_ALPHA = 0.9;
	const LEAF_NODE_HIGHLIGHT_ALPHA = 0.22;

	const RING_RADIUS_MULTIPLIER = 0.19;
	const LEAF_RING_RADIUS_MULTIPLIER = 0.12;

	const SPRING_STRENGTH = 1;
	const DAMPING = 0.1;
	const DRIFT_AMPLITUDE = 7 * SCENE_SCALE;
	const DRIFT_SPEED = 0.9;

	const ZOOM_DURATION = 650;
	const ZOOM_SCALE = 2.6;

	const LABEL_PADDING_X = 12 * SCENE_SCALE;
	const LABEL_PADDING_Y = 8 * SCENE_SCALE;
	const LABEL_COLLISION_ITERATIONS = 24;
	const TEXT_Z_INDEX = 1_000;
	const LABEL_COLOR = '--color-on-surface';

	// const MAIN_NODE_ASSET = mainNeuronImage;
	// const MAIN_NODE_ASSET_SIZE = 80 * SCENE_SCALE;
	// const LEAF_NODE_ASSET = leafNeuronImage;
	// const LEAF_NODE_ASSET_SIZE = 40 * SCENE_SCALE;
	const MAIN_NODE_ASSET = mainNeuronImage;
	const MAIN_NODE_ASSET_SIZE = 0 * SCENE_SCALE;
	const LEAF_NODE_ASSET = leafNeuronImage;
	const LEAF_NODE_ASSET_SIZE = 0 * SCENE_SCALE;

	const BUBBLE_PADDING_X = 40 * SCENE_SCALE;
	const BUBBLE_PADDING_Y = 30 * SCENE_SCALE;
	const BUBBLE_ALPHA = 0.15;

	let containerEl: HTMLDivElement;
	let canvasEl: HTMLCanvasElement;
	let app: Application | undefined;
	let clusters: Cluster[] = [];
	let world: Container | undefined;
	let resizeObserver: ResizeObserver | undefined;
	let tickerCallback: ((ticker: { deltaTime: number }) => void) | undefined;

	let focusedClusterId: number | null = null;
	let cameraTween: {
		start: number;
		duration: number;
		fromX: number;
		fromY: number;
		fromScale: number;
		toX: number;
		toY: number;
		toScale: number;
	} | null = null;

	function readCssColor(varName: string, fallback: string): number {
		if (typeof window === 'undefined') return new Color(fallback).toNumber();
		const raw = getComputedStyle(document.documentElement).getPropertyValue(varName).trim();
		try {
			return new Color(raw || fallback).toNumber();
		} catch {
			return new Color(fallback).toNumber();
		}
	}

	function easeOutCubic(t: number): number {
		return 1 - Math.pow(1 - t, 3);
	}

	type LabelBox = {
		left: number;
		right: number;
		top: number;
		bottom: number;
	};

	type CircleObstacle = {
		x: number;
		y: number;
		radius: number;
	};

	function getMainLabelBox(main: MainNode): LabelBox {
		const centerX = main.homeX;
		const topY = main.homeY + getLabelOffset(main, MAIN_NODE_CIRCLE_SIZE, 8 * SCENE_SCALE);

		return {
			left: centerX - main.labelWidth / 2 - LABEL_PADDING_X,
			right: centerX + main.labelWidth / 2 + LABEL_PADDING_X,
			top: topY - LABEL_PADDING_Y,
			bottom: topY + main.labelHeight + LABEL_PADDING_Y
		};
	}

	function getLeafLabelBox(leaf: LeafNode): LabelBox {
		const centerX = leaf.homeX;
		const topY = leaf.homeY + getLabelOffset(leaf, LEAF_NODE_CIRCLE_SIZE, 4 * SCENE_SCALE);

		return {
			left: centerX - leaf.labelWidth / 2 - LABEL_PADDING_X,
			right: centerX + leaf.labelWidth / 2 + LABEL_PADDING_X,
			top: topY - LABEL_PADDING_Y,
			bottom: topY + leaf.labelHeight + LABEL_PADDING_Y
		};
	}

	function getLabelOffset(node: MainNode | LeafNode, nodeRadius: number, gap: number): number {
		const assetRadius = node.image ? node.image.height / 2 : 0;
		return Math.max(nodeRadius, assetRadius) + gap;
	}

	function resolveLabelBoxAgainstCircle(leaf: LeafNode, circle: CircleObstacle) {
		const box = getLeafLabelBox(leaf);
		const nearestX = Math.max(box.left, Math.min(circle.x, box.right));
		const nearestY = Math.max(box.top, Math.min(circle.y, box.bottom));
		const dx = nearestX - circle.x;
		const dy = nearestY - circle.y;
		const distanceSq = dx * dx + dy * dy;
		const minDistance = circle.radius + 2 * SCENE_SCALE;

		if (distanceSq >= minDistance * minDistance) return;

		if (distanceSq === 0) {
			const boxCenterX = (box.left + box.right) / 2;
			const boxCenterY = (box.top + box.bottom) / 2;
			const pushX = boxCenterX >= circle.x ? 1 : -1;
			const pushY = boxCenterY >= circle.y ? 1 : -1;
			leaf.homeX += pushX * minDistance;
			leaf.homeY += pushY * minDistance;
			return;
		}

		const distance = Math.sqrt(distanceSq);
		const overlap = minDistance - distance;
		leaf.homeX += (dx / distance) * overlap;
		leaf.homeY += (dy / distance) * overlap;
	}

	function resolveLabelCollisions() {
		const mainLabels = clusters.map((cluster) => ({ cluster, box: getMainLabelBox(cluster.main) }));
		const leaves = clusters.flatMap((cluster) => cluster.leaves.map((leaf) => ({ cluster, leaf })));
		const nodeObstacles = clusters.flatMap((cluster) => [
			{ clusterId: cluster.id, x: cluster.main.homeX, y: cluster.main.homeY, radius: MAIN_NODE_CIRCLE_SIZE + LABEL_PADDING_Y },
			...cluster.leaves.map((leaf) => ({
				clusterId: cluster.id,
				x: leaf.homeX,
				y: leaf.homeY,
				radius: LEAF_NODE_CIRCLE_SIZE + LABEL_PADDING_Y
			}))
		]);

		for (let iteration = 0; iteration < LABEL_COLLISION_ITERATIONS; iteration += 1) {
			for (let i = 0; i < leaves.length; i += 1) {
				const current = leaves[i];

				for (let j = i + 1; j < leaves.length; j += 1) {
					const other = leaves[j];
					const currentBox = getLeafLabelBox(current.leaf);
					const otherBox = getLeafLabelBox(other.leaf);
					const overlapX = Math.min(currentBox.right, otherBox.right) - Math.max(currentBox.left, otherBox.left);
					const overlapY = Math.min(currentBox.bottom, otherBox.bottom) - Math.max(currentBox.top, otherBox.top);

					if (overlapX <= 0 || overlapY <= 0) continue;

					if (overlapX < overlapY) {
						const direction = current.leaf.homeX <= other.leaf.homeX ? -1 : 1;
						const shift = overlapX / 2 + 1;
						current.leaf.homeX += shift * direction;
						other.leaf.homeX -= shift * direction;
					} else {
						const direction = current.leaf.homeY <= other.leaf.homeY ? -1 : 1;
						const shift = overlapY / 2 + 1;
						current.leaf.homeY += shift * direction;
						other.leaf.homeY -= shift * direction;
					}
				}

				for (const obstacle of mainLabels) {
					if (obstacle.cluster.id === current.cluster.id) continue;

					const currentBox = getLeafLabelBox(current.leaf);
					const overlapX = Math.min(currentBox.right, obstacle.box.right) - Math.max(currentBox.left, obstacle.box.left);
					const overlapY = Math.min(currentBox.bottom, obstacle.box.bottom) - Math.max(currentBox.top, obstacle.box.top);

					if (overlapX <= 0 || overlapY <= 0) continue;

					if (overlapX < overlapY) {
						const direction = current.leaf.homeX <= obstacle.cluster.main.homeX ? -1 : 1;
						current.leaf.homeX += (overlapX + 1) * direction;
					} else {
						const direction = current.leaf.homeY <= obstacle.cluster.main.homeY ? -1 : 1;
						current.leaf.homeY += (overlapY + 1) * direction;
					}
				}

				for (const obstacle of nodeObstacles) {
					if (obstacle.clusterId === current.cluster.id) continue;
					resolveLabelBoxAgainstCircle(current.leaf, obstacle);
				}
			}
		}
	}

	/** Computes fixed "home" positions for all clusters so everything fits in one viewport. */
	function computeLayout(width: number, height: number) {
		const cx = width / 2;
		const cy = height / 2;
		const shortSide = Math.min(width, height);
		const overviewRingRadius = shortSide * RING_RADIUS_MULTIPLIER * SCENE_SCALE;
		const overviewLeafRingRadius = Math.max(shortSide * LEAF_RING_RADIUS_MULTIPLIER * SCENE_SCALE, 60 * SCENE_SCALE);

		if (focusedClusterId === null) {
			const count = clusters.length;

			clusters.forEach((cluster, index) => {
				const angle = (index / count) * Math.PI * 2 - Math.PI / 2;
				const mainX = cx + Math.cos(angle) * overviewRingRadius;
				const mainY = cy + Math.sin(angle) * overviewRingRadius;

				cluster.main.homeX = mainX;
				cluster.main.homeY = mainY;
				if (cluster.main.x === 0 && cluster.main.y === 0) {
					cluster.main.x = mainX;
					cluster.main.y = mainY;
				}

				const leafCount = cluster.leaves.length;
				cluster.leaves.forEach((leaf, leafIndex) => {
					const spread = Math.PI * 0.85;
					const leafAngle = angle - spread / 2 + (leafCount > 1 ? (leafIndex / (leafCount - 1)) * spread : 0);
					const leafX = mainX + Math.cos(leafAngle) * overviewLeafRingRadius;
					const leafY = mainY + Math.sin(leafAngle) * overviewLeafRingRadius;

					leaf.homeX = leafX;
					leaf.homeY = leafY;
					if (leaf.x === 0 && leaf.y === 0) {
						leaf.x = leafX;
						leaf.y = leafY;
					}
				});
			});
		} else {
			const focusedCluster = clusters.find((cluster) => cluster.id === focusedClusterId);
			const otherClusters = clusters.filter((cluster) => cluster.id !== focusedClusterId);
			const focusLeafRingRadius = Math.max(shortSide * LEAF_RING_RADIUS_MULTIPLIER * 1.4 * SCENE_SCALE, 84 * SCENE_SCALE);
			const focusMainRingRadius = Math.max(shortSide * 0.34 * SCENE_SCALE, 210 * SCENE_SCALE);

			if (focusedCluster) {
				focusedCluster.main.homeX = cx;
				focusedCluster.main.homeY = cy;

				const leafCount = focusedCluster.leaves.length;
				focusedCluster.leaves.forEach((leaf, leafIndex) => {
					const leafAngle = -Math.PI / 2 + (leafIndex / leafCount) * Math.PI * 2;
					leaf.homeX = cx + Math.cos(leafAngle) * focusLeafRingRadius;
					leaf.homeY = cy + Math.sin(leafAngle) * focusLeafRingRadius;
				});
			}

			otherClusters.forEach((cluster, index) => {
				const angle = -Math.PI / 2 + (index / otherClusters.length) * Math.PI * 2;
				const mainX = cx + Math.cos(angle) * focusMainRingRadius;
				const mainY = cy + Math.sin(angle) * focusMainRingRadius;

				cluster.main.homeX = mainX;
				cluster.main.homeY = mainY;

				cluster.leaves.forEach((leaf) => {
					leaf.homeX = mainX;
					leaf.homeY = mainY;
				});
			});
		}

		resolveLabelCollisions();
	}

	function createNodeImage(texture: Texture | undefined, size: number): Sprite | undefined {
		if (!texture) return undefined;

		const image = new Sprite(texture);
		image.anchor.set(0.5);
		const scale = size / Math.max(texture.width, texture.height);
		image.scale.set(scale);
		image.eventMode = 'none';
		return image;
	}

	async function loadNodeTexture(asset: string): Promise<Texture | undefined> {
		try {
			return await Assets.load<Texture>({
				src: asset,
				data: { resolution: 4 }
			});
		} catch {
			return undefined;
		}
	}

	function updateClusterBubble(cluster: Cluster) {
		const nodes = [cluster.main, ...cluster.leaves];
		const bounds = nodes.reduce(
			(current, node) => {
				const radius = node === cluster.main ? MAIN_NODE_CIRCLE_SIZE : LEAF_NODE_CIRCLE_SIZE;
				const labelGap = (node === cluster.main ? 8 : 4) * SCENE_SCALE;
				const labelTop = node.y + getLabelOffset(node, radius, labelGap);
				const labelLeft = node.x - node.labelWidth / 2;

				return {
					left: Math.min(current.left, node.x - radius, labelLeft),
					right: Math.max(current.right, node.x + radius, labelLeft + node.labelWidth),
					top: Math.min(current.top, node.y - radius),
					bottom: Math.max(current.bottom, node.y + radius, labelTop + node.labelHeight)
				};
			},
			{ left: Infinity, right: -Infinity, top: Infinity, bottom: -Infinity }
		);

		const centerX = (bounds.left + bounds.right) / 2;
		const centerY = (bounds.top + bounds.bottom) / 2;
		const radiusX = (bounds.right - bounds.left) / 2 + BUBBLE_PADDING_X;
		const radiusY = (bounds.bottom - bounds.top) / 2 + BUBBLE_PADDING_Y;

		cluster.bubble.clear().ellipse(centerX, centerY, radiusX, radiusY).fill({
			color: cluster.color,
			alpha: BUBBLE_ALPHA
		});
	}

	function buildClusters(
		textResolution: number,
		mainNodeTexture: Texture | undefined,
		leafNodeTexture: Texture | undefined
	): Cluster[] {
		const brainColors = [1, 2, 3, 4, 5].map((n) => readCssColor(`--color-brain-${n}`, '#c8ddf2'));
		const textColor = readCssColor(LABEL_COLOR, '#181c21');

		return sections.slice(0, 5).map((section, index) => {
			const color = brainColors[index % brainColors.length];
			const container = new Container();
			container.sortableChildren = true;
			const bubble = new Graphics();
			bubble.zIndex = 0;
			bubble.eventMode = 'none';
			const links = new Graphics();
			links.zIndex = 1;
			container.addChild(bubble, links);

			const mainGlow = new Graphics()
			.circle(0, 0, MAIN_NODE_HIGHLIGHT_RADIUS)
			.fill({ color, alpha: MAIN_NODE_HIGHLIGHT_ALPHA });
			mainGlow.filters = [new BlurFilter({ strength: 10 * SCENE_SCALE, quality: 3 })];

			const mainGraphic = new Graphics()
				.circle(0, 0, MAIN_NODE_CIRCLE_SIZE)
				.fill({ color, alpha: MAIN_NODE_ALPHA });
			mainGraphic.eventMode = 'static';
			mainGraphic.cursor = 'pointer';

			const mainLabel = new Text({
				text: section.name,
				resolution: textResolution,
				style: {
					fontFamily: 'var(--font-body)',
					fontSize: 15 * SCENE_SCALE,
					fontWeight: '600',
					fill: textColor
				}
			});
			mainLabel.anchor.set(0.5, 0);
			mainLabel.zIndex = TEXT_Z_INDEX;
			const mainImage = createNodeImage(mainNodeTexture, MAIN_NODE_ASSET_SIZE);

			container.addChild(mainGlow, mainGraphic);
			if (mainImage) container.addChild(mainImage);
			container.addChild(mainLabel);

			const main: MainNode = {
				graphic: mainGraphic,
				glow: mainGlow,
				label: mainLabel,
				imageSrc: MAIN_NODE_ASSET,
				image: mainImage,
				labelWidth: mainLabel.getLocalBounds().width,
				labelHeight: mainLabel.getLocalBounds().height,
				homeX: 0,
				homeY: 0,
				x: 0,
				y: 0,
				vx: 0,
				vy: 0,
				phaseX: Math.random() * Math.PI * 2,
				phaseY: Math.random() * Math.PI * 2
			};
			mainLabel.y = getLabelOffset(main, MAIN_NODE_CIRCLE_SIZE, 8 * SCENE_SCALE);

			const leaves: LeafNode[] = section.phrases.map((phrase) => {
				const leafGlow = new Graphics().circle(0, 0, LEAF_NODE_HIGHLIGHT_RADIUS).fill({
					color,
					alpha: LEAF_NODE_HIGHLIGHT_ALPHA
				});
				leafGlow.filters = [new BlurFilter({ strength: 5 * SCENE_SCALE, quality: 2 })];

				const leafGraphic = new Graphics()
					.circle(0, 0, LEAF_NODE_CIRCLE_SIZE)
					.fill({ color, alpha: LEAF_NODE_ALPHA });

				const label = new Text({
					text: phrase,
					resolution: textResolution,
					style: {
						fontFamily: 'var(--font-body)',
						fontSize: 11 * SCENE_SCALE,
						fill: textColor,
						align: 'center'
					}
				});
			label.anchor.set(0.5, 0);
			label.zIndex = TEXT_Z_INDEX;
			label.alpha = 0.85;
				const image = createNodeImage(leafNodeTexture, LEAF_NODE_ASSET_SIZE);

				container.addChild(leafGlow, leafGraphic);
				if (image) container.addChild(image);
				container.addChild(label);

				const leaf = {
					graphic: leafGraphic,
					glow: leafGlow,
					label,
					imageSrc: LEAF_NODE_ASSET,
					image,
					labelWidth: label.getLocalBounds().width,
					labelHeight: label.getLocalBounds().height,
					homeX: 0,
					homeY: 0,
					x: 0,
					y: 0,
					vx: 0,
					vy: 0,
					phaseX: Math.random() * Math.PI * 2,
					phaseY: Math.random() * Math.PI * 2
				};
				label.y = getLabelOffset(leaf, LEAF_NODE_CIRCLE_SIZE, 4 * SCENE_SCALE);
				return leaf;
			});

			mainGraphic.on('pointertap', (event) => {
				event.stopPropagation();
				toggleFocus(section.id);
			});

			return { id: section.id, color, container, bubble, links, main, leaves };
		});
	}

	function toggleFocus(id: number) {
		if (!app) return;

		if (focusedClusterId === id) {
			focusedClusterId = null;
		} else {
			focusedClusterId = id;
		}

		if (world) {
			world.position.set(0, 0);
			world.scale.set(1);
		}

		cameraTween = null;
		computeLayout(app.screen.width, app.screen.height);
	}

	function startCameraTween(target: { x: number; y: number; scale: number }) {
		if (!world) return;
		cameraTween = {
			start: performance.now(),
			duration: ZOOM_DURATION,
			fromX: world.x,
			fromY: world.y,
			fromScale: world.scale.x,
			toX: target.x,
			toY: target.y,
			toScale: target.scale
		};
	}

	function updateCameraTween() {
		if (!cameraTween || !world) return;
		const elapsed = performance.now() - cameraTween.start;
		const t = Math.min(elapsed / cameraTween.duration, 1);
		const eased = easeOutCubic(t);

		world.x = cameraTween.fromX + (cameraTween.toX - cameraTween.fromX) * eased;
		world.y = cameraTween.fromY + (cameraTween.toY - cameraTween.fromY) * eased;
		const scale = cameraTween.fromScale + (cameraTween.toScale - cameraTween.fromScale) * eased;
		world.scale.set(scale);

		if (t >= 1) {
			cameraTween = null;
		}
	}

	function updateDimming() {
		clusters.forEach((cluster) => {
			const isDimmed = focusedClusterId !== null && cluster.id !== focusedClusterId;
			const targetAlpha = isDimmed ? 0.12 : 1;
			cluster.container.alpha += (targetAlpha - cluster.container.alpha) * 0.15;
		});
	}

	function stepNode(node: MainNode | LeafNode, elapsedMs: number, deltaTime: number) {
		const t = elapsedMs / 1000;
		const driftX = Math.sin(t * DRIFT_SPEED + node.phaseX) * DRIFT_AMPLITUDE;
		const driftY = Math.cos(t * DRIFT_SPEED * 0.85 + node.phaseY) * DRIFT_AMPLITUDE;

		const targetX = node.homeX + driftX;
		const targetY = node.homeY + driftY;

		const ax = (targetX - node.x) * SPRING_STRENGTH;
		const ay = (targetY - node.y) * SPRING_STRENGTH;

		node.vx = (node.vx + ax) * DAMPING;
		node.vy = (node.vy + ay) * DAMPING;

		node.x += node.vx * deltaTime;
		node.y += node.vy * deltaTime;
	}

	function drawLinks(cluster: Cluster) {
		cluster.links.clear();
		cluster.leaves.forEach((leaf) => {
			cluster.links
				.moveTo(cluster.main.x, cluster.main.y)
				.lineTo(leaf.x, leaf.y)
				.stroke({ width: 1.2 * SCENE_SCALE, color: cluster.color, alpha: 0.35 });
		});
	}

	function renderFrame(elapsedMs: number, deltaTime: number) {
		clusters.forEach((cluster) => {
			stepNode(cluster.main, elapsedMs, deltaTime);
			cluster.main.graphic.position.set(cluster.main.x, cluster.main.y);
			cluster.main.glow.position.set(cluster.main.x, cluster.main.y);
			cluster.main.image?.position.set(cluster.main.x, cluster.main.y);
			cluster.main.label.position.set(
				cluster.main.x,
				cluster.main.y + getLabelOffset(cluster.main, MAIN_NODE_CIRCLE_SIZE, 8 * SCENE_SCALE)
			);

			cluster.leaves.forEach((leaf) => {
				stepNode(leaf, elapsedMs, deltaTime);
				leaf.graphic.position.set(leaf.x, leaf.y);
				leaf.glow.position.set(leaf.x, leaf.y);
				leaf.image?.position.set(leaf.x, leaf.y);
				leaf.label.position.set(leaf.x, leaf.y + getLabelOffset(leaf, LEAF_NODE_CIRCLE_SIZE, 4 * SCENE_SCALE));
			});

			updateClusterBubble(cluster);
			drawLinks(cluster);
		});

		updateCameraTween();
		updateDimming();
	}

	function handleResize() {
		if (!app) return;
		if (world) {
			world.position.set(0, 0);
			world.scale.set(1);
		}
		cameraTween = null;
		computeLayout(app.screen.width, app.screen.height);
	}

	onMount(() => {
		let cancelled = false;
		const startTime = performance.now();

		(async () => {
			const instance = new Application();
			const rendererResolution = Math.min(window.devicePixelRatio || 1, 2);
			await instance.init({
				resizeTo: containerEl,
				canvas: canvasEl,
				backgroundAlpha: 0,
				antialias: true,
				autoDensity: true,
				resolution: rendererResolution
			});

			if (cancelled) {
				instance.destroy({ releaseGlobalResources: true }, { children: true, texture: true, textureSource: true });
				return;
			}

			app = instance;

			const worldContainer = new Container();
			world = worldContainer;
			instance.stage.addChild(worldContainer);

			instance.stage.eventMode = 'static';
			instance.stage.hitArea = instance.screen;
			instance.stage.on('pointertap', () => {
				if (focusedClusterId !== null) {
					toggleFocus(focusedClusterId);
				}
			});

			const [mainNodeTexture, leafNodeTexture] = await Promise.all([
				loadNodeTexture(MAIN_NODE_ASSET),
				loadNodeTexture(LEAF_NODE_ASSET)
			]);
			if (cancelled) {
				instance.destroy({ releaseGlobalResources: true }, { children: true, texture: true, textureSource: true });
				return;
			}

			clusters = buildClusters(rendererResolution * ZOOM_SCALE, mainNodeTexture, leafNodeTexture);
			clusters.forEach((cluster) => worldContainer.addChild(cluster.container));
			computeLayout(instance.screen.width, instance.screen.height);

			tickerCallback = (ticker) => {
				renderFrame(performance.now() - startTime, ticker.deltaTime);
			};
			instance.ticker.add(tickerCallback);

			resizeObserver = new ResizeObserver(() => handleResize());
			resizeObserver.observe(containerEl);
		})();

		return () => {
			cancelled = true;
		};
	});

	onDestroy(() => {
		resizeObserver?.disconnect();
		if (app) {
			if (tickerCallback) app.ticker.remove(tickerCallback);
			app.destroy({ releaseGlobalResources: true }, { children: true, texture: true, textureSource: true });
			app = undefined;
		}
		clusters = [];
		world = undefined;
	});
</script>

<div class="mind-map" bind:this={containerEl} role="img" aria-label="Mind map of collected memories">
	<canvas bind:this={canvasEl}></canvas>
</div>

<style lang="scss">
	.mind-map {
		width: 100%;
		height: 100%;

		:global(canvas) {
			display: block;
			width: 100%;
			height: 100%;
		}
	}
</style>
