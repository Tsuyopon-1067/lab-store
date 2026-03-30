<script lang="ts">
    /**
     * バーコードリーダー入力を受け付けるコンポーネント
     * カメラ（ZXing-js）とキーボード入力（USBバーコードリーダー）の両方をサポート
     *
     * 使用方法:
     * <BarcodeInput on:scan={(e) => handleScan(e.detail)} />
     */

    import { onMount, onDestroy } from "svelte";
    import {
        BrowserMultiFormatOneDReader,
        type IScannerControls,
    } from "@zxing/browser";

    interface ScanEvent {
        barcode: string;
    }

    let input = $state("");
    let inputRef: HTMLInputElement | undefined;
    let videoRef = $state<HTMLVideoElement | undefined>(undefined);
    let cameraError = $state("");
    let cameraActive = $state(false);
    let scannerControls: IScannerControls | null = null;

    // BrowserMultiFormatReader with options
    const reader = new BrowserMultiFormatOneDReader(undefined, {
        delayBetweenScanAttempts: 500,
        delayBetweenScanSuccess: 1500, // 同一バーコードの連続誤発火防止
    });

    // バーコードスキャン完了時のイベント発火
    const dispatchScan = (barcode: string) => {
        const event = new CustomEvent<ScanEvent>("scan", {
            detail: { barcode },
        });
        window.dispatchEvent(event);
    };

    // カメラからのバーコード検出ハンドラー
    const handleCameraResult = (result: any, error: any) => {
        if (result) {
            const barcode = result.getText();
            if (barcode && barcode.trim().length > 0) {
                dispatchScan(barcode);
            }
        }
        // エラーは無視（スキャン待機中は常にエラーが出る）
    };

    // カメラスキャン開始
    async function startCamera() {
        if (!videoRef) return;

        try {
            cameraError = "";

            // カメラの設定を定義
            const constraints = {
                video: {
                    facingMode: "environment",
                    width: { ideal: 1920 },
                    height: { ideal: 1080 },
                    aspectRatio: {
                        ideal: window.innerWidth / window.innerHeight,
                    },
                    focusMode: "continuous",
                },
            };

            scannerControls = await reader.decodeFromConstraints(
                constraints,
                videoRef,
                handleCameraResult,
            );
            cameraActive = true;
        } catch (err) {
            const errMsg = err instanceof Error ? err.message : String(err);
            // カメラアクセス権限がない、または利用不可な場合
            cameraError = `カメラが利用できません: ${errMsg}。キーボード入力でご利用ください。`;
            cameraActive = false;
            // キーボード入力にフォーカス
            inputRef?.focus();
        }
    }

    // カメラスキャン停止
    function stopCamera() {
        if (scannerControls) {
            scannerControls.stop();
            scannerControls = null;
            cameraActive = false;
        }
    }

    // Enterキーでスキャン完了と判定（キーボード入力）
    function handleKeyDown(e: KeyboardEvent) {
        if (e.key === "Enter") {
            const barcode = input.trim();
            if (barcode.length > 0) {
                dispatchScan(barcode);
                input = "";
            }
        }
    }

    onMount(() => {
        // カメラを自動起動
        startCamera();
        // キーボード入力にもフォーカス
        inputRef?.focus();
    });

    onDestroy(() => {
        // コンポーネント削除時にカメラを停止
        stopCamera();
    });
</script>

<div class="barcode-input-container">
    {#if cameraError}
        <div class="error-message">
            <strong>⚠️ {cameraError}</strong>
        </div>
    {/if}

    <div
        class="camera-section"
        style="display: {cameraActive ? 'block' : 'none'}"
    >
        <video bind:this={videoRef} class="camera-preview" playsinline></video>
        <div class="camera-label">
            🎥 カメラでバーコードをスキャンしてください
        </div>
    </div>

    {#if cameraActive}
        <div class="divider">
            <span>または</span>
        </div>
    {/if}

    <input
        bind:this={inputRef}
        bind:value={input}
        type="text"
        placeholder="バーコードをスキャンしてください"
        onkeydown={handleKeyDown}
        aria-label="バーコードスキャン入力"
        class="barcode-input"
    />
</div>

<style>
    .barcode-input-container {
        width: 100%;
    }

    .error-message {
        background-color: #fff3cd;
        color: #856404;
        border: 1px solid #ffeaa7;
        border-radius: 4px;
        padding: 1rem;
        margin-bottom: 1rem;
        font-size: 0.95rem;
    }

    .camera-section {
        position: relative;
        margin-bottom: 1.5rem;
        border-radius: 8px;
        overflow: hidden;
        background-color: #000;
        aspect-ratio: 4 / 3;
    }

    .camera-preview {
        width: 100%;
        height: 100%;
        object-fit: cover;
    }

    .camera-label {
        position: absolute;
        bottom: 1rem;
        left: 50%;
        transform: translateX(-50%);
        background-color: rgba(0, 0, 0, 0.7);
        color: white;
        padding: 0.5rem 1rem;
        border-radius: 20px;
        font-size: 0.9rem;
        white-space: nowrap;
    }

    .divider {
        display: flex;
        align-items: center;
        margin: 1.5rem 0;
        color: #999;
        font-size: 0.9rem;
    }

    .divider::before,
    .divider::after {
        content: "";
        flex: 1;
        height: 1px;
        background-color: #ddd;
    }

    .divider span {
        padding: 0 1rem;
    }

    .barcode-input {
        width: 100%;
        padding: 0.75rem;
        border: 2px solid #0066cc;
        border-radius: 4px;
        font-size: 1.1rem;
        text-align: center;
        letter-spacing: 0.1em;
    }

    .barcode-input:focus {
        outline: none;
        border-color: #0052a3;
        box-shadow: 0 0 0 3px rgba(0, 102, 204, 0.15);
    }

    @media (max-width: 768px) {
        .camera-section {
            aspect-ratio: 16 / 9;
            margin-bottom: 1rem;
        }

        .camera-label {
            font-size: 0.85rem;
            padding: 0.4rem 0.8rem;
        }

        .divider {
            margin: 1rem 0;
        }

        .barcode-input {
            padding: 0.6rem;
            font-size: 1rem;
        }
    }
</style>
