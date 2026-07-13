<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { Application, Assets, Container, Graphics, Sprite, Text, Texture, Color, BlurFilter } from 'pixi.js';
	import neuron1Image from '$lib/assets/neurons/neuron_1.png';
	import neuron2Image from '$lib/assets/neurons/neuron_2.png';
	import neuron3Image from '$lib/assets/neurons/neuron_3.png';
	import neuron4Image from '$lib/assets/neurons/neuron_4.png';

	type MindMapSection = {
		id: number;
		name: string;
		anchor: {
			x: number;
			y: number;
		};
		rotationSpeed: number;
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
		rotationSpeed: number;
		scale: number;
		main: MainNode;
		leaves: LeafNode[];
	};

	let sections: MindMapSection[] = [];

	const SCENE_SCALE = 1.4;

	const MAIN_NODE_CIRCLE_SIZE = 10 * SCENE_SCALE;
	const MAIN_NODE_HIGHLIGHT_RADIUS = 60 * SCENE_SCALE;
	const MAIN_NODE_ALPHA = 1;
	const MAIN_NODE_HIGHLIGHT_ALPHA = 0.2;
	const NEURON_ASSETS = [neuron1Image, neuron2Image, neuron3Image, neuron4Image];
	const MAIN_NODE_ASSET_SIZE = 200 * SCENE_SCALE;
	const MAIN_NODE_ASSET_ALPHA = 1;

	const LEAF_NODE_CIRCLE_SIZE = 7 * SCENE_SCALE;
	const LEAF_NODE_HIGHLIGHT_RADIUS = 0 * SCENE_SCALE;
	const LEAF_NODE_ALPHA = 1;
	const LEAF_NODE_HIGHLIGHT_ALPHA = 0.8;
	const LEAF_RING_RADIUS_MULTIPLIER = 0.1;
	const LEAF_MIN_DISTANCE = 56 * SCENE_SCALE;
	const LEAF_MIN_GAP_RATIO = 0.72;

	const BUBBLE_PADDING_X = 90 * SCENE_SCALE;
	const BUBBLE_PADDING_Y = 60 * SCENE_SCALE;
	const BUBBLE_ALPHA = 0.5;
	const BLOB_PROFILES = [
		[1, 0.96, 0.92, 1.04, 1.08, 1.02, 0.94, 0.9, 0.98, 1.06, 1.04, 0.96, 1.02, 0.93, 1.07, 1.01, 0.95, 1.03, 0.92, 0.98],
		[0.95, 1.02, 1.08, 1.04, 0.96, 0.91, 0.94, 1.03, 1.06, 1.1, 1.02, 0.98, 0.94, 1.02, 0.9, 0.96, 1.04, 1.01, 0.92, 1.05],
		[1.06, 1.02, 0.95, 0.9, 0.96, 1.04, 1.08, 1.03, 0.97, 0.92, 1.01, 1.07, 1.1, 1.04, 0.94, 0.91, 1.01, 1.06, 0.98, 0.96],
		[0.92, 0.98, 1.04, 1.1, 1.06, 0.96, 0.92, 1.02, 1.08, 1.03, 0.9, 0.95, 1.05, 1.07, 0.98, 0.93, 1.03, 1.06, 0.94, 1.01],
		[1.04, 1.01, 0.9, 0.94, 1.07, 1.09, 1.03, 0.98, 0.93, 1.05, 1.01, 0.91, 1.06, 1.08, 0.97, 0.94, 1.02, 1.05, 0.96, 1.03]
	] as const;

	const NODE_FILL_COLOR = 0xE5E5E5;
	const LINK_COLOR = 0xE5E5E5;
	const LINK_WIDTH = 2 * SCENE_SCALE;

	const SPRING_STRENGTH = 1;
	const DAMPING = 0.1;
	const DRIFT_AMPLITUDE = 7 * SCENE_SCALE;
	const DRIFT_SPEED = 0.9;

	const CLUSTER_ZOOM_SCALE = 1.8;
	const CLUSTER_BACKGROUND_SCALE = 0.65;
	const CLUSTER_ZOOM_EASING = 0.12;
	const TEXT_RESOLUTION_SCALE = 2.6;

	const LABEL_PADDING_X = 12 * SCENE_SCALE;
	const LABEL_PADDING_Y = 8 * SCENE_SCALE;
	const MAIN_LABEL_GAP = -70 * SCENE_SCALE;
	const LABEL_COLLISION_ITERATIONS = 24;
	const TEXT_Z_INDEX = 1_000;
	const LABEL_COLOR = '--color-on-surface';

	let containerEl: HTMLDivElement;
	let canvasEl: HTMLCanvasElement;
	let app: Application | undefined;
	let clusters: Cluster[] = [];
	let world: Container | undefined;
	let resizeObserver: ResizeObserver | undefined;
	let tickerCallback: ((ticker: { deltaTime: number }) => void) | undefined;

	let focusedClusterId: number | null = null;

	function readCssColor(varName: string, fallback: string): number {
		if (typeof window === 'undefined') return new Color(fallback).toNumber();
		const raw = getComputedStyle(document.documentElement).getPropertyValue(varName).trim();
		try {
			return new Color(raw || fallback).toNumber();
		} catch {
			return new Color(fallback).toNumber();
		}
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
		const topY = main.homeY + getLabelOffset(main, MAIN_NODE_CIRCLE_SIZE, MAIN_LABEL_GAP);

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

	function getLeafRingRadius(baseRadius: number, count: number, minimumAngularGap: number): number {
		if (count <= 1) return baseRadius;

		const radiusForMinimumDistance = LEAF_MIN_DISTANCE / (2 * Math.sin(minimumAngularGap / 2));

		return Math.max(baseRadius, radiusForMinimumDistance);
	}

	function getCircularLeafAngles(count: number): number[] {
		if (count <= 1) return [0];

		const equalSpacing = (Math.PI * 2) / count;
		const maximumOffset = equalSpacing * 0.1;

		return Array.from({ length: count }, (_, index) => {
			const equalAngle = index * equalSpacing;
			const offset = Math.sin((index + 1) * 1.7) * maximumOffset;
			return equalAngle + offset;
		});
	}

	/** Computes fixed "home" positions for all clusters so everything fits in one viewport. */
	function computeLayout(width: number, height: number) {
		const cx = width / 2;
		const cy = height / 2;
		const shortSide = Math.min(width, height);
		const overviewLeafRingRadius = Math.max(shortSide * LEAF_RING_RADIUS_MULTIPLIER * SCENE_SCALE, 60 * SCENE_SCALE);

		if (focusedClusterId === null) {
			clusters.forEach((cluster) => {
				const section = sections.find((candidate) => candidate.id === cluster.id);
				if (!section) return;

				const mainX = width * section.anchor.x;
				const mainY = height * section.anchor.y;
				const angle = Math.atan2(mainY - cy, mainX - cx);

				cluster.main.homeX = mainX;
				cluster.main.homeY = mainY;
				if (cluster.main.x === 0 && cluster.main.y === 0) {
					cluster.main.x = mainX;
					cluster.main.y = mainY;
				}

				const leafCount = cluster.leaves.length;
				const overviewMaximumGap = leafCount > 1 ? (Math.PI * 2) / leafCount : 0;
				const overviewMinimumGap = overviewMaximumGap * LEAF_MIN_GAP_RATIO;
				const leafRingRadius = getLeafRingRadius(overviewLeafRingRadius, leafCount, overviewMinimumGap);
				const leafAngles = getCircularLeafAngles(leafCount);
				cluster.leaves.forEach((leaf, leafIndex) => {
					const leafAngle = angle - Math.PI / 2 + (leafAngles[leafIndex] ?? 0);
					const leafX = mainX + Math.cos(leafAngle) * leafRingRadius;
					const leafY = mainY + Math.sin(leafAngle) * leafRingRadius;

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
			const focusLeafRingRadius = Math.max(shortSide * LEAF_RING_RADIUS_MULTIPLIER * 1.4 * SCENE_SCALE, 84 * SCENE_SCALE);

			if (focusedCluster) {
				const focusedSection = sections.find((section) => section.id === focusedCluster.id);
				const focusedAnchor = focusedSection
					? { x: width * focusedSection.anchor.x, y: height * focusedSection.anchor.y }
					: { x: focusedCluster.main.homeX, y: focusedCluster.main.homeY };
				const centerClusterOverview = clusters
					.slice()
					.sort((first, second) => {
						const firstDistance = Math.hypot(first.main.homeX - cx, first.main.homeY - cy);
						const secondDistance = Math.hypot(second.main.homeX - cx, second.main.homeY - cy);
						return firstDistance - secondDistance;
					})[0];
				const isCenterClusterSelected = centerClusterOverview?.id === focusedCluster.id;
				const centerSwapCluster = clusters
					.filter((cluster) => cluster.id !== focusedCluster.id)
					.sort((first, second) => {
						const firstDistance = Math.hypot(first.main.homeX - cx, first.main.homeY - cy);
						const secondDistance = Math.hypot(second.main.homeX - cx, second.main.homeY - cy);
						return firstDistance - secondDistance;
					})[0];

				if (!isCenterClusterSelected && centerSwapCluster) {
					const offsetX = focusedAnchor.x - centerSwapCluster.main.homeX;
					const offsetY = focusedAnchor.y - centerSwapCluster.main.homeY;
					centerSwapCluster.main.homeX += offsetX;
					centerSwapCluster.main.homeY += offsetY;
					centerSwapCluster.leaves.forEach((leaf) => {
						leaf.homeX += offsetX;
						leaf.homeY += offsetY;
					});
				}

				focusedCluster.main.homeX = cx;
				focusedCluster.main.homeY = cy;

				const leafCount = focusedCluster.leaves.length;
				const focusMinimumGap = ((Math.PI * 2) / leafCount) * (1 - 0.2);
				const leafRingRadius = getLeafRingRadius(focusLeafRingRadius, leafCount, focusMinimumGap);
				const leafAngles = getCircularLeafAngles(leafCount);
				focusedCluster.leaves.forEach((leaf, leafIndex) => {
					const leafAngle = -Math.PI / 2 + (leafAngles[leafIndex] ?? 0);
					leaf.homeX = cx + Math.cos(leafAngle) * leafRingRadius;
					leaf.homeY = cy + Math.sin(leafAngle) * leafRingRadius;
				});
			}
		}

		if (focusedClusterId === null) resolveLabelCollisions();
	}

	function createNodeImage(texture: Texture | undefined, size: number, alpha: number): Sprite | undefined {
		if (!texture) return undefined;

		const image = new Sprite(texture);
		image.anchor.set(0.5);
		const scale = size / Math.max(texture.width, texture.height);
		image.scale.set(scale);
		image.alpha = alpha;
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

	function drawBlob(
		graphic: Graphics,
		centerX: number,
		centerY: number,
		radiusX: number,
		radiusY: number,
		profileIndex: number,
		color: number
	) {
		const profile = BLOB_PROFILES[profileIndex % BLOB_PROFILES.length];
		const points = profile.map((radius, index) => {
			const angle = -Math.PI / 2 + (index / profile.length) * Math.PI * 2;
			return {
				x: centerX + Math.cos(angle) * radiusX * radius,
				y: centerY + Math.sin(angle) * radiusY * radius
			};
		});

		graphic.clear().moveTo(points[0].x, points[0].y);
		points.forEach((point, index) => {
			const previous = points[(index - 1 + points.length) % points.length];
			const next = points[(index + 1) % points.length];
			const following = points[(index + 2) % points.length];
			const incomingTangent = { x: next.x - previous.x, y: next.y - previous.y };
			const outgoingTangent = { x: following.x - point.x, y: following.y - point.y };
			const nextPoint = next;
			const handle = 0.18;

			graphic.bezierCurveTo(
				point.x + incomingTangent.x * handle,
				point.y + incomingTangent.y * handle,
				nextPoint.x - outgoingTangent.x * handle,
				nextPoint.y - outgoingTangent.y * handle,
				nextPoint.x,
				nextPoint.y
			);
		});

		graphic.closePath().fill({ color, alpha: BUBBLE_ALPHA });
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

		drawBlob(cluster.bubble, centerX, centerY, radiusX, radiusY, cluster.id, cluster.color);
	}

	function buildClusters(
		textResolution: number,
		neuronTextures: (Texture | undefined)[]
	): Cluster[] {
		const brainColors = [1, 2, 3, 4, 5].map((n) => readCssColor(`--color-brain-${n}`, '#c8ddf2'));
		const textColor = readCssColor(LABEL_COLOR, '#181c21');

		return sections.slice(0, 5).map((section, index) => {
			const color = brainColors[index % brainColors.length];
			const container = new Container();
			container.sortableChildren = true;
			container.eventMode = 'static';
			container.cursor = 'pointer';
			const bubble = new Graphics();
			bubble.zIndex = 0;
			bubble.eventMode = 'static';
			bubble.cursor = 'pointer';
			const links = new Graphics();
			links.zIndex = 1;
			container.addChild(bubble, links);

			const mainGlow = new Graphics()
				.circle(0, 0, MAIN_NODE_HIGHLIGHT_RADIUS)
				.fill({ color, alpha: MAIN_NODE_HIGHLIGHT_ALPHA });

			const mainGraphic = new Graphics()
				.circle(0, 0, MAIN_NODE_CIRCLE_SIZE)
			.fill({ color: NODE_FILL_COLOR, alpha: MAIN_NODE_ALPHA });
			mainGraphic.eventMode = 'none';

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
			const neuronAssetIndex = index % NEURON_ASSETS.length;
			const mainImage = createNodeImage(
				neuronTextures[neuronAssetIndex],
				MAIN_NODE_ASSET_SIZE,
				MAIN_NODE_ASSET_ALPHA
			);

			container.addChild(mainGlow, mainGraphic);
			if (mainImage) container.addChild(mainImage);
			container.addChild(mainLabel);

			const main: MainNode = {
				graphic: mainGraphic,
				glow: mainGlow,
				label: mainLabel,
				imageSrc: NEURON_ASSETS[neuronAssetIndex],
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
			mainLabel.y = getLabelOffset(main, MAIN_NODE_CIRCLE_SIZE, MAIN_LABEL_GAP);

			const leaves: LeafNode[] = section.phrases.map((phrase) => {
				const leafGlow = new Graphics().circle(0, 0, LEAF_NODE_HIGHLIGHT_RADIUS).fill({
					color,
					alpha: LEAF_NODE_HIGHLIGHT_ALPHA
				});

				const leafGraphic = new Graphics()
					.circle(0, 0, LEAF_NODE_CIRCLE_SIZE)
					.fill({ color: NODE_FILL_COLOR, alpha: LEAF_NODE_ALPHA });

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
				container.addChild(leafGlow, leafGraphic);
				container.addChild(label);

				const leaf = {
					graphic: leafGraphic,
					glow: leafGlow,
					label,
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

			container.on('pointertap', (event) => {
				event.stopPropagation();
				toggleFocus(section.id);
			});

			return {
				id: section.id,
				color,
				container,
				bubble,
				links,
				rotationSpeed: section.rotationSpeed,
				scale: 1,
				main,
				leaves
			};
		});
	}

	function toggleFocus(id: number) {
		if (!app) return;

		if (focusedClusterId === id) {
			focusedClusterId = null;
		} else {
			focusedClusterId = id;
		}

		computeLayout(app.screen.width, app.screen.height);
	}

	function updateDimming() {
		clusters.forEach((cluster) => {
			const isDimmed = focusedClusterId !== null && cluster.id !== focusedClusterId;
			const targetAlpha = isDimmed ? 0.12 : 1;
			cluster.container.alpha += (targetAlpha - cluster.container.alpha) * 0.15;
		});
	}

	function stepNode(
		node: MainNode | LeafNode,
		elapsedMs: number,
		deltaTime: number,
		rotationCenter?: { x: number; y: number },
		rotation = 0
	) {
		const t = elapsedMs / 1000;
		const driftX = Math.sin(t * DRIFT_SPEED + node.phaseX) * DRIFT_AMPLITUDE;
		const driftY = Math.cos(t * DRIFT_SPEED * 0.85 + node.phaseY) * DRIFT_AMPLITUDE;
		let homeX = node.homeX;
		let homeY = node.homeY;

		if (rotationCenter) {
			const relativeX = node.homeX - rotationCenter.x;
			const relativeY = node.homeY - rotationCenter.y;
			const cos = Math.cos(rotation);
			const sin = Math.sin(rotation);
			homeX = rotationCenter.x + relativeX * cos - relativeY * sin;
			homeY = rotationCenter.y + relativeX * sin + relativeY * cos;
		}

		const targetX = homeX + driftX;
		const targetY = homeY + driftY;

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
				.stroke({ width: LINK_WIDTH, color: LINK_COLOR, alpha: 0.35 });
		});
	}

	function renderFrame(elapsedMs: number, deltaTime: number) {
		clusters.forEach((cluster) => {
			const rotation = (elapsedMs / 1000) * cluster.rotationSpeed;
			const targetScale =
				focusedClusterId === null
					? 1
					: focusedClusterId === cluster.id
						? CLUSTER_ZOOM_SCALE
						: CLUSTER_BACKGROUND_SCALE;
			cluster.scale += (targetScale - cluster.scale) * CLUSTER_ZOOM_EASING;
			stepNode(cluster.main, elapsedMs, deltaTime);
			cluster.container.pivot.set(cluster.main.x, cluster.main.y);
			cluster.container.position.set(cluster.main.x, cluster.main.y);
			cluster.container.scale.set(cluster.scale);
			cluster.main.graphic.position.set(cluster.main.x, cluster.main.y);
			cluster.main.glow.position.set(cluster.main.x, cluster.main.y);
			if (cluster.main.image) {
				cluster.main.image.position.set(cluster.main.x, cluster.main.y);
				cluster.main.image.rotation = rotation;
			}
			cluster.main.label.position.set(
				cluster.main.x,
				cluster.main.y + getLabelOffset(cluster.main, MAIN_NODE_CIRCLE_SIZE, MAIN_LABEL_GAP)
			);

			cluster.leaves.forEach((leaf) => {
				stepNode(leaf, elapsedMs, deltaTime, cluster.main, rotation);
				leaf.graphic.position.set(leaf.x, leaf.y);
				leaf.glow.position.set(leaf.x, leaf.y);
				if (leaf.image) {
					leaf.image.position.set(leaf.x, leaf.y);
					leaf.image.rotation = rotation;
				}
				leaf.label.position.set(leaf.x, leaf.y + getLabelOffset(leaf, LEAF_NODE_CIRCLE_SIZE, 4 * SCENE_SCALE));
			});

			updateClusterBubble(cluster);
			drawLinks(cluster);
		});

		updateDimming();
	}

	function handleResize() {
		if (!app) return;
		computeLayout(app.screen.width, app.screen.height);
	}

		onMount(() => {
		let cancelled = false;
		const startTime = performance.now();

		(async () => {
			const dataResponse = await fetch('/mindmap.json');
			if (!dataResponse.ok) return;
			sections = (await dataResponse.json()) as MindMapSection[];

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

			const neuronTextures = await Promise.all(NEURON_ASSETS.map((asset) => loadNodeTexture(asset)));
			if (cancelled) {
				instance.destroy({ releaseGlobalResources: true }, { children: true, texture: true, textureSource: true });
				return;
			}

			clusters = buildClusters(rendererResolution * TEXT_RESOLUTION_SCALE, neuronTextures);
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
