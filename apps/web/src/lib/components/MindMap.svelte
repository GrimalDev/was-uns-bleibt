<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { SvelteMap } from 'svelte/reactivity';
	import { Application, Assets, Container, Graphics, Sprite, Text, Texture, Color, BlurFilter } from 'pixi.js';
	import neuron1Image from '$lib/assets/neurons/neuron_1.png';
	import neuron2Image from '$lib/assets/neurons/neuron_2.png';
	import neuron3Image from '$lib/assets/neurons/neuron_3.png';
	import neuron4Image from '$lib/assets/neurons/neuron_4.png';
	import neuron5Image from '$lib/assets/neurons/neuron_5.png';
	import bubble1Image from '$lib/assets/bubbles/bubble_1.png';
	import bubble2Image from '$lib/assets/bubbles/bubble_2.png';
	import bubble3Image from '$lib/assets/bubbles/bubble_3.png';
	import bubble4Image from '$lib/assets/bubbles/bubble_4.png';
	import bubble5Image from '$lib/assets/bubbles/bubble_5.png';

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

	type Answer = {
		id: number;
		brain_part_id: number;
		phrase: string;
	};

	type AnswerEvent = {
		type: 'snapshot' | 'answer';
		answers?: Answer[];
		answer?: Answer;
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
		labelHighlight: Graphics;
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
		bubble: Sprite;
		links: Graphics;
		rotationSpeed: number;
		scale: number;
		main: MainNode;
		leaves: LeafNode[];
	};

	let sections: MindMapSection[] = [];

	// Scene and asset configuration
	const SCENE_SCALE_FACTOR = 1.4;
	const NEURON_TEXTURE_ASSETS = [neuron1Image, neuron2Image, neuron3Image, neuron4Image, neuron5Image];
	const BUBBLE_TEXTURE_ASSETS = [bubble1Image, bubble2Image, bubble3Image, bubble4Image, bubble5Image];

	// Main node appearance
	const MAIN_NODE_RADIUS = 10 * SCENE_SCALE_FACTOR;
	const MAIN_NODE_GLOW_RADIUS_X = 60 * SCENE_SCALE_FACTOR;
	const MAIN_NODE_GLOW_RADIUS_Y = 60 * SCENE_SCALE_FACTOR;
	const MAIN_NODE_OPACITY = 1;
	const MAIN_NODE_GLOW_OPACITY = 0.2;
	const MAIN_NODE_IMAGE_SIZE = 200 * SCENE_SCALE_FACTOR;
	const MAIN_NODE_IMAGE_OPACITY = 1;

	// Leaf node appearance and layout
	const LEAF_NODE_RADIUS = 7 * SCENE_SCALE_FACTOR;
	const LEAF_NODE_GLOW_RADIUS = 0 * SCENE_SCALE_FACTOR;
	const LEAF_NODE_OPACITY = 1;
	const LEAF_NODE_GLOW_OPACITY = 0.8;
	const LEAF_RING_RADIUS_SCALE = 0.1;
	const CENTER_NODE_LEAF_DISTANCE_SCALE = 0.95;
	const MINIMUM_LEAF_DISTANCE = 56 * SCENE_SCALE_FACTOR;
	const MINIMUM_LEAF_ANGULAR_GAP_RATIO = 0.72;

	// Bubble image layout and shape profiles
	const BUBBLE_PADDING_HORIZONTAL = 90 * SCENE_SCALE_FACTOR;
	const BUBBLE_PADDING_VERTICAL = 60 * SCENE_SCALE_FACTOR;
	const BUBBLE_SCALE = 0.95;
	const BUBBLE_OPACITY = 0.5;
	const BUBBLE_IMAGE_CONTENT_WIDTH_RATIO = 780 / 1920;
	const BUBBLE_IMAGE_CONTENT_HEIGHT_RATIO = 490 / 1080;
	const BUBBLE_LABEL_HORIZONTAL_OFFSETS = [-0.03, 0.035, -0.03, 0, 0] as const;
	const BLOB_SHAPE_PROFILES = [
		[1, 0.96, 0.92, 1.04, 1.08, 1.02, 0.94, 0.9, 0.98, 1.06, 1.04, 0.96, 1.02, 0.93, 1.07, 1.01, 0.95, 1.03, 0.92, 0.98],
		[0.95, 1.02, 1.08, 1.04, 0.96, 0.91, 0.94, 1.03, 1.06, 1.1, 1.02, 0.98, 0.94, 1.02, 0.9, 0.96, 1.04, 1.01, 0.92, 1.05],
		[1.06, 1.02, 0.95, 0.9, 0.96, 1.04, 1.08, 1.03, 0.97, 0.92, 1.01, 1.07, 1.1, 1.04, 0.94, 0.91, 1.01, 1.06, 0.98, 0.96],
		[0.92, 0.98, 1.04, 1.1, 1.06, 0.96, 0.92, 1.02, 1.08, 1.03, 0.9, 0.95, 1.05, 1.07, 0.98, 0.93, 1.03, 1.06, 0.94, 1.01],
		[1.04, 1.01, 0.9, 0.94, 1.07, 1.09, 1.03, 0.98, 0.93, 1.05, 1.01, 0.91, 1.06, 1.08, 0.97, 0.94, 1.02, 1.05, 0.96, 1.03]
	] as const;

	// Node and connection rendering
	const NODE_FILL_COLOR_HEX = 0xE5E5E5;
	const CONNECTION_LINK_COLOR = 0xE5E5E5;
	const CONNECTION_LINK_WIDTH = 2 * SCENE_SCALE_FACTOR;

	// Node motion
	const NODE_SPRING_STRENGTH = 1;
	const MOTION_DAMPING = 0.1;
	const NODE_DRIFT_AMPLITUDE = 7 * SCENE_SCALE_FACTOR;
	const NODE_DRIFT_SPEED = 0.9;

	// Cluster focus and rendering
	const FOCUSED_CLUSTER_SCALE = 1.8;
	const UNFOCUSED_CLUSTER_SCALE = 0.65;
	const CLUSTER_SCALE_EASING = 0.12;
	const LABEL_RENDER_RESOLUTION_SCALE = 2.6;

	// Label layout and styling
	const LABEL_PADDING_HORIZONTAL = 20 * SCENE_SCALE_FACTOR;
	const LABEL_PADDING_VERTICAL = 4 * SCENE_SCALE_FACTOR;
	const MAIN_LABEL_VERTICAL_GAP = -70 * SCENE_SCALE_FACTOR;
	const MAX_LABEL_COLLISION_ITERATIONS = 24;
	const LABEL_Z_INDEX = 1_000;
	const MAIN_LABEL_CSS_COLOR = '--color-on-surface';
	const LEAF_LABEL_CSS_COLOR = '--color-inverse-on-surface';

	let containerEl: HTMLDivElement;
	let canvasEl: HTMLCanvasElement;
	let app: Application | undefined;
	let clusters: Cluster[] = [];
	let world: Container | undefined;
	let resizeObserver: ResizeObserver | undefined;
	let tickerCallback: ((ticker: { deltaTime: number }) => void) | undefined;
	let answerSocket: WebSocket | undefined;
	let textResolution = 0;

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
		const topY = main.homeY + getLabelOffset(main, MAIN_NODE_RADIUS, MAIN_LABEL_VERTICAL_GAP);

		return {
			left: centerX - main.labelWidth / 2 - LABEL_PADDING_HORIZONTAL,
			right: centerX + main.labelWidth / 2 + LABEL_PADDING_HORIZONTAL,
			top: topY - LABEL_PADDING_VERTICAL,
			bottom: topY + main.labelHeight + LABEL_PADDING_VERTICAL
		};
	}

	function getLeafLabelBox(leaf: LeafNode): LabelBox {
		const centerX = leaf.homeX;
		const topY = leaf.homeY + getLabelOffset(leaf, LEAF_NODE_RADIUS, 4 * SCENE_SCALE_FACTOR);

		return {
			left: centerX - leaf.labelWidth / 2 - LABEL_PADDING_HORIZONTAL,
			right: centerX + leaf.labelWidth / 2 + LABEL_PADDING_HORIZONTAL,
			top: topY - LABEL_PADDING_VERTICAL,
			bottom: topY + leaf.labelHeight + LABEL_PADDING_VERTICAL
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
		const minDistance = circle.radius + 2 * SCENE_SCALE_FACTOR;

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
					{ clusterId: cluster.id, x: cluster.main.homeX, y: cluster.main.homeY, radius: MAIN_NODE_RADIUS + LABEL_PADDING_VERTICAL },
			...cluster.leaves.map((leaf) => ({
				clusterId: cluster.id,
				x: leaf.homeX,
				y: leaf.homeY,
				radius: LEAF_NODE_RADIUS + LABEL_PADDING_VERTICAL
			}))
		]);

		for (let iteration = 0; iteration < MAX_LABEL_COLLISION_ITERATIONS; iteration += 1) {
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

		const radiusForMinimumDistance = MINIMUM_LEAF_DISTANCE / (2 * Math.sin(minimumAngularGap / 2));

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
		const overviewLeafRingRadius = Math.max(shortSide * LEAF_RING_RADIUS_SCALE * SCENE_SCALE_FACTOR, 60 * SCENE_SCALE_FACTOR);

		if (focusedClusterId === null) {
			const centerCluster = clusters
				.slice()
				.sort((first, second) => {
					const firstSection = sections.find((section) => section.id === first.id);
					const secondSection = sections.find((section) => section.id === second.id);
					const firstDistance = firstSection
						? Math.hypot(width * firstSection.anchor.x - cx, height * firstSection.anchor.y - cy)
						: Infinity;
					const secondDistance = secondSection
						? Math.hypot(width * secondSection.anchor.x - cx, height * secondSection.anchor.y - cy)
						: Infinity;
					return firstDistance - secondDistance;
				})[0];

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
				const overviewMinimumGap = overviewMaximumGap * MINIMUM_LEAF_ANGULAR_GAP_RATIO;
				const leafRingRadius =
					getLeafRingRadius(overviewLeafRingRadius, leafCount, overviewMinimumGap) *
					(cluster.id === centerCluster?.id ? CENTER_NODE_LEAF_DISTANCE_SCALE : 1);
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
			const focusLeafRingRadius = Math.max(
				shortSide * LEAF_RING_RADIUS_SCALE * 1.4 * SCENE_SCALE_FACTOR,
				84 * SCENE_SCALE_FACTOR
			);

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
			const leafRingRadius =
				getLeafRingRadius(focusLeafRingRadius, leafCount, focusMinimumGap) *
				CENTER_NODE_LEAF_DISTANCE_SCALE;
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
		color: number,
		alpha = BUBBLE_OPACITY
	) {
		const profile = BLOB_SHAPE_PROFILES[profileIndex % BLOB_SHAPE_PROFILES.length];
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

		graphic.closePath().fill({ color, alpha });
	}

	function updateClusterBubble(cluster: Cluster) {
		const nodes = [cluster.main, ...cluster.leaves];
		const bounds = nodes.reduce(
			(current, node) => {
			const radius = node === cluster.main ? MAIN_NODE_RADIUS : LEAF_NODE_RADIUS;
			const labelGap = (node === cluster.main ? 8 : 4) * SCENE_SCALE_FACTOR;
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
		const radiusX = (bounds.right - bounds.left) / 2 + BUBBLE_PADDING_HORIZONTAL;
		const radiusY = (bounds.bottom - bounds.top) / 2 + BUBBLE_PADDING_VERTICAL;

		cluster.bubble.position.set(centerX, centerY);
		cluster.bubble.width = ((radiusX * 2) / BUBBLE_IMAGE_CONTENT_WIDTH_RATIO) * BUBBLE_SCALE;
		cluster.bubble.height = ((radiusY * 2) / BUBBLE_IMAGE_CONTENT_HEIGHT_RATIO) * BUBBLE_SCALE;
	}

	function updateMainLabelPosition(cluster: Cluster) {
		const bubbleBounds = cluster.bubble.getLocalBounds();
		const bubbleCenterX = cluster.bubble.x + bubbleBounds.x + bubbleBounds.width / 2;
		const labelX = bubbleCenterX + cluster.bubble.width * (BUBBLE_LABEL_HORIZONTAL_OFFSETS[cluster.id - 1] ?? 0);
		const labelY = cluster.bubble.y + cluster.bubble.height / 5;

		cluster.main.label.position.set(labelX, labelY);
		cluster.main.labelHighlight.position.set(labelX, labelY);
	}

	function createLeafNode(
		phrase: string,
		textResolution: number,
		color: number,
		leafLabelColor: number,
		container: Container
	): LeafNode {
		const leafGlow = new Graphics().circle(0, 0, LEAF_NODE_GLOW_RADIUS).fill({
			color,
			alpha: LEAF_NODE_GLOW_OPACITY
		});
		const leafGraphic = new Graphics()
			.circle(0, 0, LEAF_NODE_RADIUS)
			.fill({ color: NODE_FILL_COLOR_HEX, alpha: LEAF_NODE_OPACITY });
		const label = new Text({
			text: phrase,
			resolution: textResolution,
			style: {
				fontFamily: 'var(--font-body)',
				fontSize: 12 * SCENE_SCALE_FACTOR,
				fontWeight: '700',
				fill: leafLabelColor,
				align: 'center'
			}
		});
		label.anchor.set(0.5, 0);
		label.zIndex = LABEL_Z_INDEX;
		label.alpha = 0.85;
		container.addChild(leafGlow, leafGraphic, label);

		const leaf: LeafNode = {
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
		label.y = getLabelOffset(leaf, LEAF_NODE_RADIUS, 4 * SCENE_SCALE_FACTOR);
		return leaf;
	}

	function syncAnswer(answer: Answer) {
		const section = sections.find((candidate) => candidate.id === answer.brain_part_id);
		const cluster = clusters.find((candidate) => candidate.id === answer.brain_part_id);
		if (!section || !cluster || section.phrases.includes(answer.phrase)) return;

		section.phrases = [...section.phrases, answer.phrase];
		cluster.leaves.push(
			createLeafNode(
				answer.phrase,
				textResolution,
				cluster.color,
				readCssColor(LEAF_LABEL_CSS_COLOR, '#404751'),
				cluster.container
			)
		);
	}

	function syncAnswers(answers: Answer[]) {
		for (const answer of answers) syncAnswer(answer);
		if (app) computeLayout(app.screen.width, app.screen.height);
	}

	function connectAnswerSocket() {
		const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
		answerSocket = new WebSocket(`${protocol}//${window.location.host}/api/answers/ws`);
		answerSocket.addEventListener('message', (message) => {
			if (typeof message.data !== 'string') return;

			try {
				const event = JSON.parse(message.data) as AnswerEvent;
				if (event.type === 'snapshot') syncAnswers(event.answers ?? []);
				if (event.type === 'answer' && event.answer) syncAnswer(event.answer);
				if (event.type === 'answer' && app) computeLayout(app.screen.width, app.screen.height);
			} catch {
				return;
			}
		});
	}

	function buildClusters(
		textResolution: number,
		neuronTextures: (Texture | undefined)[],
		bubbleTextures: (Texture | undefined)[]
	): Cluster[] {
		const brainColors = [1, 2, 3, 4, 5].map((n) => readCssColor(`--color-brain-${n}`, '#c8ddf2'));
		const mainLabelColor = readCssColor(MAIN_LABEL_CSS_COLOR, '#181c21');
		const leafLabelColor = readCssColor(LEAF_LABEL_CSS_COLOR, '#404751');

		return sections.slice(0, 5).map((section, index) => {
			const color = brainColors[index % brainColors.length];
			const container = new Container();
			container.sortableChildren = true;
			container.eventMode = 'static';
			container.cursor = 'pointer';
			const bubble = new Sprite(bubbleTextures[index] ?? Texture.EMPTY);
			bubble.zIndex = 0;
			bubble.anchor.set(0.5);
			bubble.eventMode = 'none';
			const links = new Graphics();
			links.zIndex = 1;
			container.addChild(bubble, links);

			const mainGlow = new Graphics();
			drawBlob(
				mainGlow,
				0,
				0,
				MAIN_NODE_GLOW_RADIUS_X,
				MAIN_NODE_GLOW_RADIUS_Y,
				section.id,
				color,
				MAIN_NODE_GLOW_OPACITY
			);

			const mainGraphic = new Graphics()
				.circle(0, 0, MAIN_NODE_RADIUS)
				.fill({ color: NODE_FILL_COLOR_HEX, alpha: MAIN_NODE_OPACITY });
			mainGraphic.eventMode = 'none';

			const mainLabel = new Text({
					text: section.name.split('/').map((part) => part.trim()).join('\n'),
				resolution: textResolution,
					style: {
						fontFamily: 'var(--font-body)',
					fontSize: 15 * SCENE_SCALE_FACTOR,
						fontWeight: '600',
						fill: mainLabelColor,
						align: 'center'
					}
		});
		mainLabel.anchor.set(0.5, 0);
			mainLabel.zIndex = LABEL_Z_INDEX;
		const mainLabelHighlight = new Graphics()
			.roundRect(
				-mainLabel.getLocalBounds().width / 2 - LABEL_PADDING_HORIZONTAL,
				-LABEL_PADDING_VERTICAL,
				mainLabel.getLocalBounds().width + LABEL_PADDING_HORIZONTAL * 2,
				mainLabel.getLocalBounds().height + LABEL_PADDING_VERTICAL * 2,
				LABEL_PADDING_VERTICAL
			)
			.fill({ color, alpha: 1 });
			mainLabelHighlight.zIndex = LABEL_Z_INDEX - 1;
			const neuronAssetIndex = section.id - 1;
			const mainImage = createNodeImage(
				neuronTextures[neuronAssetIndex],
				MAIN_NODE_IMAGE_SIZE,
				MAIN_NODE_IMAGE_OPACITY
			);

			container.addChild(mainGlow, mainGraphic, mainLabelHighlight);
			if (mainImage) container.addChild(mainImage);
			container.addChild(mainLabel);

			const main: MainNode = {
				graphic: mainGraphic,
				glow: mainGlow,
				label: mainLabel,
				labelHighlight: mainLabelHighlight,
				imageSrc: NEURON_TEXTURE_ASSETS[neuronAssetIndex],
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
			mainLabel.y = getLabelOffset(main, MAIN_NODE_RADIUS, MAIN_LABEL_VERTICAL_GAP);

			const leaves = section.phrases.map((phrase) =>
				createLeafNode(phrase, textResolution, color, leafLabelColor, container)
			);

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
		const driftX = Math.sin(t * NODE_DRIFT_SPEED + node.phaseX) * NODE_DRIFT_AMPLITUDE;
		const driftY = Math.cos(t * NODE_DRIFT_SPEED * 0.85 + node.phaseY) * NODE_DRIFT_AMPLITUDE;
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

		const ax = (targetX - node.x) * NODE_SPRING_STRENGTH;
		const ay = (targetY - node.y) * NODE_SPRING_STRENGTH;

		node.vx = (node.vx + ax) * MOTION_DAMPING;
		node.vy = (node.vy + ay) * MOTION_DAMPING;

		node.x += node.vx * deltaTime;
		node.y += node.vy * deltaTime;
	}

	function drawLinks(cluster: Cluster) {
		cluster.links.clear();
		cluster.leaves.forEach((leaf) => {
			cluster.links
				.moveTo(cluster.main.x, cluster.main.y)
				.lineTo(leaf.x, leaf.y)
				.stroke({ width: CONNECTION_LINK_WIDTH, color: CONNECTION_LINK_COLOR, alpha: 0.35 });
		});
	}

	function renderFrame(elapsedMs: number, deltaTime: number) {
		clusters.forEach((cluster) => {
			const rotation = (elapsedMs / 1000) * cluster.rotationSpeed;
			const targetScale =
				focusedClusterId === null
					? 1
					: focusedClusterId === cluster.id
					? FOCUSED_CLUSTER_SCALE
					: UNFOCUSED_CLUSTER_SCALE;
			cluster.scale += (targetScale - cluster.scale) * CLUSTER_SCALE_EASING;
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

			cluster.leaves.forEach((leaf) => {
				stepNode(leaf, elapsedMs, deltaTime, cluster.main, rotation);
				leaf.graphic.position.set(leaf.x, leaf.y);
				leaf.glow.position.set(leaf.x, leaf.y);
				if (leaf.image) {
					leaf.image.position.set(leaf.x, leaf.y);
					leaf.image.rotation = rotation;
				}
				leaf.label.position.set(leaf.x, leaf.y + getLabelOffset(leaf, LEAF_NODE_RADIUS, 4 * SCENE_SCALE_FACTOR));
			});

			updateClusterBubble(cluster);
			updateMainLabelPosition(cluster);
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

			const mindMapSections = (await dataResponse.json()) as Omit<MindMapSection, 'phrases'>[];
			let answers: Answer[] = [];
			try {
				const answersResponse = await fetch('/api/answers');
				if (answersResponse.ok) answers = (await answersResponse.json()) as Answer[];
			} catch {
				answers = [];
			}

			const phrasesByPartId = new SvelteMap<number, string[]>();
			for (const answer of answers) {
				const phrases = phrasesByPartId.get(answer.brain_part_id) ?? [];
				phrases.push(answer.phrase);
				phrasesByPartId.set(answer.brain_part_id, phrases);
			}
			sections = mindMapSections.map((section) => ({
				...section,
				phrases: phrasesByPartId.get(section.id) ?? []
			}));

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

			const neuronTextures = await Promise.all(NEURON_TEXTURE_ASSETS.map((asset) => loadNodeTexture(asset)));
			const bubbleTextures = await Promise.all(BUBBLE_TEXTURE_ASSETS.map((asset) => loadNodeTexture(asset)));
			if (cancelled) {
				instance.destroy({ releaseGlobalResources: true }, { children: true, texture: true, textureSource: true });
				return;
			}

			clusters = buildClusters(rendererResolution * LABEL_RENDER_RESOLUTION_SCALE, neuronTextures, bubbleTextures);
			textResolution = rendererResolution * LABEL_RENDER_RESOLUTION_SCALE;
			clusters.forEach((cluster) => worldContainer.addChild(cluster.container));
			computeLayout(instance.screen.width, instance.screen.height);
			connectAnswerSocket();

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
		answerSocket?.close();
		answerSocket = undefined;
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
