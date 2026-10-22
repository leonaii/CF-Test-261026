
CF上号器_v1.17.exe:	file format coff-x86-64

Disassembly of section .text:

0000000140001000 <.text>:
14008a400: 48 ba fe ff ff ff ff ff ff 7f       	movabsq	$0x7ffffffffffffffe, %rdx # imm = 0x7FFFFFFFFFFFFFFE
14008a40a: 48 29 c2                    	subq	%rax, %rdx
14008a40d: 4c 39 c2                    	cmpq	%r8, %rdx
14008a410: 0f 82 4d df 04 00           	jb	0x1400d8363 <.text+0xd7363>
14008a416: 48 89 ca                    	movq	%rcx, %rdx
14008a419: 48 8d 8c 24 a0 00 00 00     	leaq	0xa0(%rsp), %rcx
14008a421: e8 aa 84 03 00              	callq	0x1400c28d0 <.text+0xc18d0>
14008a426: 4c 8d ac 24 10 01 00 00     	leaq	0x110(%rsp), %r13
14008a42e: 48 8d 50 10                 	leaq	0x10(%rax), %rdx
14008a432: 4c 89 ac 24 00 01 00 00     	movq	%r13, 0x100(%rsp)
14008a43a: 48 8b 08                    	movq	(%rax), %rcx
14008a43d: 48 39 d1                    	cmpq	%rdx, %rcx
14008a440: 0f 84 aa 0c 00 00           	je	0x14008b0f0 <.text+0x8a0f0>
14008a446: 48 89 8c 24 00 01 00 00     	movq	%rcx, 0x100(%rsp)
14008a44e: 48 8b 48 10                 	movq	0x10(%rax), %rcx
14008a452: 48 89 8c 24 10 01 00 00     	movq	%rcx, 0x110(%rsp)
14008a45a: 48 8b 48 08                 	movq	0x8(%rax), %rcx
14008a45e: 48 89 8c 24 08 01 00 00     	movq	%rcx, 0x108(%rsp)
14008a466: c6 40 10 00                 	movb	$0x0, 0x10(%rax)
14008a46a: 48 c7 40 08 00 00 00 00     	movq	$0x0, 0x8(%rax)
14008a472: 48 89 10                    	movq	%rdx, (%rax)
14008a475: 48 b8 0e 00 00 00 00 00 00 80       	movabsq	$-0x7ffffffffffffff2, %rax # imm = 0x800000000000000E
14008a47f: 48 03 84 24 08 01 00 00     	addq	0x108(%rsp), %rax
14008a487: 48 83 f8 0c                 	cmpq	$0xc, %rax
14008a48b: 0f 86 b6 de 04 00           	jbe	0x1400d8347 <.text+0xd7347>
14008a491: 4c 8d bc 24 00 01 00 00     	leaq	0x100(%rsp), %r15
14008a499: 41 b8 0d 00 00 00           	movl	$0xd, %r8d
14008a49f: 48 8d 15 49 08 07 00        	leaq	0x70849(%rip), %rdx     # 0x1400facef
14008a4a6: 4c 89 f9                    	movq	%r15, %rcx
14008a4a9: e8 22 84 03 00              	callq	0x1400c28d0 <.text+0xc18d0>
14008a4ae: 48 8d ac 24 30 01 00 00     	leaq	0x130(%rsp), %rbp
14008a4b6: 48 8d 48 10                 	leaq	0x10(%rax), %rcx
14008a4ba: 48 89 ac 24 20 01 00 00     	movq	%rbp, 0x120(%rsp)
14008a4c2: 48 8b 10                    	movq	(%rax), %rdx
14008a4c5: 48 39 ca                    	cmpq	%rcx, %rdx
14008a4c8: 0f 84 e2 0a 00 00           	je	0x14008afb0 <.text+0x89fb0>
14008a4ce: 48 89 94 24 20 01 00 00     	movq	%rdx, 0x120(%rsp)
14008a4d6: 48 8b 50 10                 	movq	0x10(%rax), %rdx
14008a4da: 48 89 94 24 30 01 00 00     	movq	%rdx, 0x130(%rsp)
14008a4e2: 48 8b 50 08                 	movq	0x8(%rax), %rdx
14008a4e6: 48 89 08                    	movq	%rcx, (%rax)
14008a4e9: 4c 8b 8c 24 20 01 00 00     	movq	0x120(%rsp), %r9
14008a4f1: 48 c7 40 08 00 00 00 00     	movq	$0x0, 0x8(%rax)
14008a4f9: 4c 8b 84 24 68 01 00 00     	movq	0x168(%rsp), %r8
14008a501: c6 40 10 00                 	movb	$0x0, 0x10(%rax)
14008a505: 48 8b 8c 24 60 01 00 00     	movq	0x160(%rsp), %rcx
14008a50d: 48 89 94 24 28 01 00 00     	movq	%rdx, 0x128(%rsp)
14008a515: 4a 8d 04 02                 	leaq	(%rdx,%r8), %rax
14008a519: 49 39 e9                    	cmpq	%rbp, %r9
14008a51c: 0f 84 6e 0a 00 00           	je	0x14008af90 <.text+0x89f90>
14008a522: 48 39 84 24 30 01 00 00     	cmpq	%rax, 0x130(%rsp)
14008a52a: 73 1b                       	jae	0x14008a547 <.text+0x89547>
14008a52c: 4c 8d 94 24 70 01 00 00     	leaq	0x170(%rsp), %r10
14008a534: 4c 39 d1                    	cmpq	%r10, %rcx
14008a537: 74 0e                       	je	0x14008a547 <.text+0x89547>
14008a539: 48 39 84 24 70 01 00 00     	cmpq	%rax, 0x170(%rsp)
14008a541: 0f 83 d9 07 00 00           	jae	0x14008ad20 <.text+0x89d20>
14008a547: 48 b8 fe ff ff ff ff ff ff 7f       	movabsq	$0x7ffffffffffffffe, %rax # imm = 0x7FFFFFFFFFFFFFFE
14008a551: 48 29 d0                    	subq	%rdx, %rax
14008a554: 4c 39 c0                    	cmpq	%r8, %rax
14008a557: 0f 82 d6 dd 04 00           	jb	0x1400d8333 <.text+0xd7333>
14008a55d: 48 8d 9c 24 20 01 00 00     	leaq	0x120(%rsp), %rbx
14008a565: 48 89 ca                    	movq	%rcx, %rdx
14008a568: 48 89 d9                    	movq	%rbx, %rcx
14008a56b: e8 60 83 03 00              	callq	0x1400c28d0 <.text+0xc18d0>
14008a570: 48 8d bc 24 90 01 00 00     	leaq	0x190(%rsp), %rdi
14008a578: 48 8d 50 10                 	leaq	0x10(%rax), %rdx
14008a57c: 48 89 bc 24 80 01 00 00     	movq	%rdi, 0x180(%rsp)
14008a584: 48 8b 08                    	movq	(%rax), %rcx
14008a587: 48 89 7c 24 40              	movq	%rdi, 0x40(%rsp)
14008a58c: 48 39 d1                    	cmpq	%rdx, %rcx
14008a58f: 0f 84 f3 0a 00 00           	je	0x14008b088 <.text+0x8a088>
14008a595: 48 89 8c 24 80 01 00 00     	movq	%rcx, 0x180(%rsp)
14008a59d: 48 8b 48 10                 	movq	0x10(%rax), %rcx
14008a5a1: 48 89 8c 24 90 01 00 00     	movq	%rcx, 0x190(%rsp)
14008a5a9: 48 8b 48 08                 	movq	0x8(%rax), %rcx
14008a5ad: 48 89 8c 24 88 01 00 00     	movq	%rcx, 0x188(%rsp)
14008a5b5: c6 40 10 00                 	movb	$0x0, 0x10(%rax)
14008a5b9: 48 c7 40 08 00 00 00 00     	movq	$0x0, 0x8(%rax)
14008a5c1: 48 89 10                    	movq	%rdx, (%rax)
14008a5c4: 48 b8 0f 00 00 00 00 00 00 80       	movabsq	$-0x7ffffffffffffff1, %rax # imm = 0x800000000000000F
14008a5ce: 48 03 84 24 88 01 00 00     	addq	0x188(%rsp), %rax
14008a5d6: 48 83 f8 0d                 	cmpq	$0xd, %rax
14008a5da: 0f 86 32 dd 04 00           	jbe	0x1400d8312 <.text+0xd7312>
14008a5e0: 48 8d 84 24 80 01 00 00     	leaq	0x180(%rsp), %rax
14008a5e8: 41 b8 0e 00 00 00           	movl	$0xe, %r8d
14008a5ee: 48 8d 15 08 07 07 00        	leaq	0x70708(%rip), %rdx     # 0x1400facfd
14008a5f5: 48 89 c1                    	movq	%rax, %rcx
14008a5f8: 48 89 44 24 58              	movq	%rax, 0x58(%rsp)
14008a5fd: e8 ce 82 03 00              	callq	0x1400c28d0 <.text+0xc18d0>
14008a602: 4c 8d a4 24 b0 01 00 00     	leaq	0x1b0(%rsp), %r12
14008a60a: 48 8d 48 10                 	leaq	0x10(%rax), %rcx
14008a60e: 4c 89 a4 24 a0 01 00 00     	movq	%r12, 0x1a0(%rsp)
14008a616: 48 8b 10                    	movq	(%rax), %rdx
14008a619: 48 39 ca                    	cmpq	%rcx, %rdx
14008a61c: 0f 84 fe 09 00 00           	je	0x14008b020 <.text+0x8a020>
14008a622: 48 89 94 24 a0 01 00 00     	movq	%rdx, 0x1a0(%rsp)
14008a62a: 48 8b 50 10                 	movq	0x10(%rax), %rdx
14008a62e: 48 89 94 24 b0 01 00 00     	movq	%rdx, 0x1b0(%rsp)
14008a636: 48 8b 50 08                 	movq	0x8(%rax), %rdx
14008a63a: 48 89 08                    	movq	%rcx, (%rax)
14008a63d: 4c 8b 8c 24 a0 01 00 00     	movq	0x1a0(%rsp), %r9
14008a645: 48 c7 40 08 00 00 00 00     	movq	$0x0, 0x8(%rax)
14008a64d: 4c 8b 84 24 e8 01 00 00     	movq	0x1e8(%rsp), %r8
14008a655: c6 40 10 00                 	movb	$0x0, 0x10(%rax)
14008a659: 48 8b 8c 24 e0 01 00 00     	movq	0x1e0(%rsp), %rcx
14008a661: 48 89 94 24 a8 01 00 00     	movq	%rdx, 0x1a8(%rsp)
14008a669: 4a 8d 04 02                 	leaq	(%rdx,%r8), %rax
14008a66d: 4d 39 e1                    	cmpq	%r12, %r9
14008a670: 0f 84 8a 09 00 00           	je	0x14008b000 <.text+0x8a000>
14008a676: 48 39 84 24 b0 01 00 00     	cmpq	%rax, 0x1b0(%rsp)
14008a67e: 73 1b                       	jae	0x14008a69b <.text+0x8969b>
14008a680: 4c 8d 94 24 f0 01 00 00     	leaq	0x1f0(%rsp), %r10
14008a688: 4c 39 d1                    	cmpq	%r10, %rcx
14008a68b: 74 0e                       	je	0x14008a69b <.text+0x8969b>
14008a68d: 48 39 84 24 f0 01 00 00     	cmpq	%rax, 0x1f0(%rsp)
14008a695: 0f 83 f5 07 00 00           	jae	0x14008ae90 <.text+0x89e90>
14008a69b: 48 b8 fe ff ff ff ff ff ff 7f       	movabsq	$0x7ffffffffffffffe, %rax # imm = 0x7FFFFFFFFFFFFFFE
14008a6a5: 48 29 d0                    	subq	%rdx, %rax
14008a6a8: 4c 39 c0                    	cmpq	%r8, %rax
14008a6ab: 0f 82 15 dc 04 00           	jb	0x1400d82c6 <.text+0xd72c6>
14008a6b1: 48 8d 9c 24 a0 01 00 00     	leaq	0x1a0(%rsp), %rbx
14008a6b9: 48 89 ca                    	movq	%rcx, %rdx
14008a6bc: 48 89 d9                    	movq	%rbx, %rcx
14008a6bf: e8 0c 82 03 00              	callq	0x1400c28d0 <.text+0xc18d0>
14008a6c4: 48 8d 9c 24 10 02 00 00     	leaq	0x210(%rsp), %rbx
14008a6cc: 48 8d 50 10                 	leaq	0x10(%rax), %rdx
14008a6d0: 48 89 9c 24 00 02 00 00     	movq	%rbx, 0x200(%rsp)
14008a6d8: 48 8b 08                    	movq	(%rax), %rcx
14008a6db: 48 89 5c 24 48              	movq	%rbx, 0x48(%rsp)
14008a6e0: 48 39 d1                    	cmpq	%rdx, %rcx
14008a6e3: 0f 84 67 09 00 00           	je	0x14008b050 <.text+0x8a050>
14008a6e9: 48 89 8c 24 00 02 00 00     	movq	%rcx, 0x200(%rsp)
14008a6f1: 48 8b 48 10                 	movq	0x10(%rax), %rcx
14008a6f5: 48 89 8c 24 10 02 00 00     	movq	%rcx, 0x210(%rsp)
14008a6fd: 48 8b 48 08                 	movq	0x8(%rax), %rcx
14008a701: 48 89 8c 24 08 02 00 00     	movq	%rcx, 0x208(%rsp)
14008a709: c6 40 10 00                 	movb	$0x0, 0x10(%rax)
14008a70d: 48 c7 40 08 00 00 00 00     	movq	$0x0, 0x8(%rax)
14008a715: 48 89 10                    	movq	%rdx, (%rax)
14008a718: 48 b8 09 00 00 00 00 00 00 80       	movabsq	$-0x7ffffffffffffff7, %rax # imm = 0x8000000000000009
14008a722: 48 03 84 24 08 02 00 00     	addq	0x208(%rsp), %rax
14008a72a: 48 83 f8 07                 	cmpq	$0x7, %rax
14008a72e: 0f 86 79 db 04 00           	jbe	0x1400d82ad <.text+0xd72ad>
14008a734: 48 8d 84 24 00 02 00 00     	leaq	0x200(%rsp), %rax
14008a73c: 41 b8 08 00 00 00           	movl	$0x8, %r8d
14008a742: 48 8d 15 c3 05 07 00        	leaq	0x705c3(%rip), %rdx     # 0x1400fad0c
14008a749: 48 89 c1                    	movq	%rax, %rcx
14008a74c: 48 89 44 24 50              	movq	%rax, 0x50(%rsp)
14008a751: e8 7a 81 03 00              	callq	0x1400c28d0 <.text+0xc18d0>
14008a756: 48 8d 9c 24 30 02 00 00     	leaq	0x230(%rsp), %rbx
14008a75e: 48 8d 50 10                 	leaq	0x10(%rax), %rdx
14008a762: 48 89 9c 24 20 02 00 00     	movq	%rbx, 0x220(%rsp)
14008a76a: 48 8b 08                    	movq	(%rax), %rcx
14008a76d: 48 39 d1                    	cmpq	%rdx, %rcx
14008a770: 0f 84 ea 07 00 00           	je	0x14008af60 <.text+0x89f60>
14008a776: 48 89 8c 24 20 02 00 00     	movq	%rcx, 0x220(%rsp)
14008a77e: 48 8b 48 10                 	movq	0x10(%rax), %rcx
14008a782: 48 89 8c 24 30 02 00 00     	movq	%rcx, 0x230(%rsp)
14008a78a: 48 8b 48 08                 	movq	0x8(%rax), %rcx
14008a78e: 48 89 10                    	movq	%rdx, (%rax)
14008a791: 48 89 8c 24 28 02 00 00     	movq	%rcx, 0x228(%rsp)
14008a799: 48 c7 40 08 00 00 00 00     	movq	$0x0, 0x8(%rax)
14008a7a1: c6 40 10 00                 	movb	$0x0, 0x10(%rax)
14008a7a5: 48 8b 84 24 28 02 00 00     	movq	0x228(%rsp), %rax
14008a7ad: 4c 8b 84 24 48 02 00 00     	movq	0x248(%rsp), %r8
14008a7b5: 4c 8b 8c 24 20 02 00 00     	movq	0x220(%rsp), %r9
14008a7bd: 48 8b 8c 24 40 02 00 00     	movq	0x240(%rsp), %rcx
14008a7c5: 4a 8d 14 00                 	leaq	(%rax,%r8), %rdx
14008a7c9: 49 39 d9                    	cmpq	%rbx, %r9
14008a7cc: 0f 84 6e 07 00 00           	je	0x14008af40 <.text+0x89f40>
14008a7d2: 48 39 94 24 30 02 00 00     	cmpq	%rdx, 0x230(%rsp)
14008a7da: 73 13                       	jae	0x14008a7ef <.text+0x897ef>
14008a7dc: 48 39 f1                    	cmpq	%rsi, %rcx
14008a7df: 74 0e                       	je	0x14008a7ef <.text+0x897ef>
14008a7e1: 48 39 94 24 50 02 00 00     	cmpq	%rdx, 0x250(%rsp)
14008a7e9: 0f 83 19 06 00 00           	jae	0x14008ae08 <.text+0x89e08>
14008a7ef: 48 ba fe ff ff ff ff ff ff 7f       	movabsq	$0x7ffffffffffffffe, %rdx # imm = 0x7FFFFFFFFFFFFFFE
14008a7f9: 48 29 c2                    	subq	%rax, %rdx
14008a7fc: 4c 39 c2                    	cmpq	%r8, %rdx
14008a7ff: 0f 82 8a da 04 00           	jb	0x1400d828f <.text+0xd728f>
14008a805: 48 8d 84 24 20 02 00 00     	leaq	0x220(%rsp), %rax
14008a80d: 48 89 ca                    	movq	%rcx, %rdx
14008a810: 48 89 c1                    	movq	%rax, %rcx
14008a813: 48 89 44 24 30              	movq	%rax, 0x30(%rsp)
14008a818: e8 b3 80 03 00              	callq	0x1400c28d0 <.text+0xc18d0>
14008a81d: 48 8d bc 24 70 02 00 00     	leaq	0x270(%rsp), %rdi
14008a825: 48 8d 50 10                 	leaq	0x10(%rax), %rdx
14008a829: 48 89 bc 24 60 02 00 00     	movq	%rdi, 0x260(%rsp)
14008a831: 48 8b 08                    	movq	(%rax), %rcx
14008a834: 48 39 d1                    	cmpq	%rdx, %rcx
14008a837: 0f 84 83 08 00 00           	je	0x14008b0c0 <.text+0x8a0c0>
14008a83d: 48 89 8c 24 60 02 00 00     	movq	%rcx, 0x260(%rsp)
14008a845: 48 8b 48 10                 	movq	0x10(%rax), %rcx
14008a849: 48 89 8c 24 70 02 00 00     	movq	%rcx, 0x270(%rsp)
14008a851: 48 8b 48 08                 	movq	0x8(%rax), %rcx
14008a855: 48 89 8c 24 68 02 00 00     	movq	%rcx, 0x268(%rsp)
14008a85d: 48 89 10                    	movq	%rdx, (%rax)
14008a860: c6 40 10 00                 	movb	$0x0, 0x10(%rax)
14008a864: 48 c7 40 08 00 00 00 00     	movq	$0x0, 0x8(%rax)
14008a86c: 48 b8 03 00 00 00 00 00 00 80       	movabsq	$-0x7ffffffffffffffd, %rax # imm = 0x8000000000000003
14008a876: 48 03 84 24 68 02 00 00     	addq	0x268(%rsp), %rax
14008a87e: 48 83 f8 01                 	cmpq	$0x1, %rax
14008a882: 0f 86 93 d9 04 00           	jbe	0x1400d821b <.text+0xd721b>
14008a888: 48 8d 84 24 60 02 00 00     	leaq	0x260(%rsp), %rax
14008a890: 41 b8 02 00 00 00           	movl	$0x2, %r8d
14008a896: 48 8d 15 78 04 07 00        	leaq	0x70478(%rip), %rdx     # 0x1400fad15
14008a89d: 48 89 c1                    	movq	%rax, %rcx
14008a8a0: 48 89 44 24 38              	movq	%rax, 0x38(%rsp)
14008a8a5: e8 26 80 03 00              	callq	0x1400c28d0 <.text+0xc18d0>
14008a8aa: 4c 8d bc 24 90 00 00 00     	leaq	0x90(%rsp), %r15
14008a8b2: 48 8d 50 10                 	leaq	0x10(%rax), %rdx
14008a8b6: 4c 89 bc 24 80 00 00 00     	movq	%r15, 0x80(%rsp)
14008a8be: 48 8b 08                    	movq	(%rax), %rcx
14008a8c1: 48 39 d1                    	cmpq	%rdx, %rcx
14008a8c4: 0f 84 46 06 00 00           	je	0x14008af10 <.text+0x89f10>
14008a8ca: 48 89 8c 24 80 00 00 00     	movq	%rcx, 0x80(%rsp)
14008a8d2: 48 8b 48 10                 	movq	0x10(%rax), %rcx
14008a8d6: 48 89 8c 24 90 00 00 00     	movq	%rcx, 0x90(%rsp)
14008a8de: 48 8b 48 08                 	movq	0x8(%rax), %rcx
14008a8e2: 48 89 8c 24 88 00 00 00     	movq	%rcx, 0x88(%rsp)
14008a8ea: 48 89 10                    	movq	%rdx, (%rax)
14008a8ed: c6 40 10 00                 	movb	$0x0, 0x10(%rax)
14008a8f1: 48 8b 8c 24 60 02 00 00     	movq	0x260(%rsp), %rcx
14008a8f9: 48 c7 40 08 00 00 00 00     	movq	$0x0, 0x8(%rax)
14008a901: 48 39 f9                    	cmpq	%rdi, %rcx
14008a904: 74 11                       	je	0x14008a917 <.text+0x89917>
14008a906: 48 8b 84 24 70 02 00 00     	movq	0x270(%rsp), %rax
14008a90e: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008a912: e8 b9 f7 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008a917: 48 8b 8c 24 20 02 00 00     	movq	0x220(%rsp), %rcx
14008a91f: 48 39 d9                    	cmpq	%rbx, %rcx
14008a922: 74 11                       	je	0x14008a935 <.text+0x89935>
14008a924: 48 8b 84 24 30 02 00 00     	movq	0x230(%rsp), %rax
14008a92c: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008a930: e8 9b f7 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008a935: 48 8b 8c 24 00 02 00 00     	movq	0x200(%rsp), %rcx
14008a93d: 48 3b 4c 24 48              	cmpq	0x48(%rsp), %rcx
14008a942: 74 11                       	je	0x14008a955 <.text+0x89955>
14008a944: 48 8b 84 24 10 02 00 00     	movq	0x210(%rsp), %rax
14008a94c: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008a950: e8 7b f7 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008a955: 48 8b 8c 24 a0 01 00 00     	movq	0x1a0(%rsp), %rcx
14008a95d: 4c 39 e1                    	cmpq	%r12, %rcx
14008a960: 74 11                       	je	0x14008a973 <.text+0x89973>
14008a962: 48 8b 84 24 b0 01 00 00     	movq	0x1b0(%rsp), %rax
14008a96a: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008a96e: e8 5d f7 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008a973: 48 8b 8c 24 80 01 00 00     	movq	0x180(%rsp), %rcx
14008a97b: 48 3b 4c 24 40              	cmpq	0x40(%rsp), %rcx
14008a980: 74 11                       	je	0x14008a993 <.text+0x89993>
14008a982: 48 8b 84 24 90 01 00 00     	movq	0x190(%rsp), %rax
14008a98a: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008a98e: e8 3d f7 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008a993: 48 8b 8c 24 20 01 00 00     	movq	0x120(%rsp), %rcx
14008a99b: 48 39 e9                    	cmpq	%rbp, %rcx
14008a99e: 74 11                       	je	0x14008a9b1 <.text+0x899b1>
14008a9a0: 48 8b 84 24 30 01 00 00     	movq	0x130(%rsp), %rax
14008a9a8: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008a9ac: e8 1f f7 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008a9b1: 48 8b 8c 24 00 01 00 00     	movq	0x100(%rsp), %rcx
14008a9b9: 4c 39 e9                    	cmpq	%r13, %rcx
14008a9bc: 74 11                       	je	0x14008a9cf <.text+0x899cf>
14008a9be: 48 8b 84 24 10 01 00 00     	movq	0x110(%rsp), %rax
14008a9c6: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008a9ca: e8 01 f7 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008a9cf: 48 8b 8c 24 a0 00 00 00     	movq	0xa0(%rsp), %rcx
14008a9d7: 4c 39 f1                    	cmpq	%r14, %rcx
14008a9da: 74 11                       	je	0x14008a9ed <.text+0x899ed>
14008a9dc: 48 8b 84 24 b0 00 00 00     	movq	0xb0(%rsp), %rax
14008a9e4: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008a9e8: e8 e3 f6 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008a9ed: 48 8b 8c 24 e0 00 00 00     	movq	0xe0(%rsp), %rcx
14008a9f5: 48 8d 84 24 f0 00 00 00     	leaq	0xf0(%rsp), %rax
14008a9fd: 48 39 c1                    	cmpq	%rax, %rcx
14008aa00: 74 11                       	je	0x14008aa13 <.text+0x89a13>
14008aa02: 48 8b 84 24 f0 00 00 00     	movq	0xf0(%rsp), %rax
14008aa0a: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008aa0e: e8 bd f6 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008aa13: 48 8b 8c 24 c0 00 00 00     	movq	0xc0(%rsp), %rcx
14008aa1b: 48 8d 84 24 d0 00 00 00     	leaq	0xd0(%rsp), %rax
14008aa23: 48 39 c1                    	cmpq	%rax, %rcx
14008aa26: 74 11                       	je	0x14008aa39 <.text+0x89a39>
14008aa28: 48 8b 84 24 d0 00 00 00     	movq	0xd0(%rsp), %rax
14008aa30: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008aa34: e8 97 f6 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008aa39: 48 8b 8c 24 60 01 00 00     	movq	0x160(%rsp), %rcx
14008aa41: 48 8d 84 24 70 01 00 00     	leaq	0x170(%rsp), %rax
14008aa49: 48 39 c1                    	cmpq	%rax, %rcx
14008aa4c: 74 11                       	je	0x14008aa5f <.text+0x89a5f>
14008aa4e: 48 8b 84 24 70 01 00 00     	movq	0x170(%rsp), %rax
14008aa56: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008aa5a: e8 71 f6 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008aa5f: 48 8b 8c 24 40 01 00 00     	movq	0x140(%rsp), %rcx
14008aa67: 48 8d 84 24 50 01 00 00     	leaq	0x150(%rsp), %rax
14008aa6f: 48 39 c1                    	cmpq	%rax, %rcx
14008aa72: 74 11                       	je	0x14008aa85 <.text+0x89a85>
14008aa74: 48 8b 84 24 50 01 00 00     	movq	0x150(%rsp), %rax
14008aa7c: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008aa80: e8 4b f6 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008aa85: 48 8b 8c 24 e0 01 00 00     	movq	0x1e0(%rsp), %rcx
14008aa8d: 48 8d 84 24 f0 01 00 00     	leaq	0x1f0(%rsp), %rax
14008aa95: 48 39 c1                    	cmpq	%rax, %rcx
14008aa98: 74 11                       	je	0x14008aaab <.text+0x89aab>
14008aa9a: 48 8b 84 24 f0 01 00 00     	movq	0x1f0(%rsp), %rax
14008aaa2: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008aaa6: e8 25 f6 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008aaab: 48 8b 8c 24 c0 01 00 00     	movq	0x1c0(%rsp), %rcx
14008aab3: 48 8d 84 24 d0 01 00 00     	leaq	0x1d0(%rsp), %rax
14008aabb: 48 39 c1                    	cmpq	%rax, %rcx
14008aabe: 74 11                       	je	0x14008aad1 <.text+0x89ad1>
14008aac0: 48 8b 84 24 d0 01 00 00     	movq	0x1d0(%rsp), %rax
14008aac8: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008aacc: e8 ff f5 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008aad1: 48 8b 8c 24 40 02 00 00     	movq	0x240(%rsp), %rcx
14008aad9: 48 39 f1                    	cmpq	%rsi, %rcx
14008aadc: 74 11                       	je	0x14008aaef <.text+0x89aef>
14008aade: 48 8b 84 24 50 02 00 00     	movq	0x250(%rsp), %rax
14008aae6: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008aaea: e8 e1 f5 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008aaef: 45 31 d2                    	xorl	%r10d, %r10d
14008aaf2: 48 8d 0d 1f 02 07 00        	leaq	0x7021f(%rip), %rcx     # 0x1400fad18
14008aaf9: 48 89 9c 24 20 02 00 00     	movq	%rbx, 0x220(%rsp)
14008ab01: c6 84 24 30 02 00 00 00     	movb	$0x0, 0x230(%rsp)
14008ab09: 48 c7 84 24 28 02 00 00 00 00 00 00 	movq	$0x0, 0x228(%rsp)
14008ab15: 48 89 b4 24 40 02 00 00     	movq	%rsi, 0x240(%rsp)
14008ab1d: 66 44 89 94 24 50 02 00 00  	movw	%r10w, 0x250(%rsp)
14008ab26: 48 c7 84 24 48 02 00 00 00 00 00 00 	movq	$0x0, 0x248(%rsp)
14008ab32: 48 89 bc 24 60 02 00 00     	movq	%rdi, 0x260(%rsp)
14008ab3a: e8 41 ab 01 00              	callq	0x1400a5680 <.text+0xa4680>
14008ab3f: 4c 8d 24 00                 	leaq	(%rax,%rax), %r12
14008ab43: 4c 89 e2                    	movq	%r12, %rdx
14008ab46: 48 d1 fa                    	sarq	%rdx
14008ab49: 48 89 94 24 00 02 00 00     	movq	%rdx, 0x200(%rsp)
14008ab51: 49 83 fc 0e                 	cmpq	$0xe, %r12
14008ab55: 0f 87 35 01 00 00           	ja	0x14008ac90 <.text+0x89c90>
14008ab5b: 48 8b 8c 24 60 02 00 00     	movq	0x260(%rsp), %rcx
14008ab63: 48 83 fa 01                 	cmpq	$0x1, %rdx
14008ab67: 0f 84 0b 01 00 00           	je	0x14008ac78 <.text+0x89c78>
14008ab6d: 31 c0                       	xorl	%eax, %eax
14008ab6f: 48 85 d2                    	testq	%rdx, %rdx
14008ab72: 0f 85 4d 01 00 00           	jne	0x14008acc5 <.text+0x89cc5>
14008ab78: 45 31 c0                    	xorl	%r8d, %r8d
14008ab7b: 48 8d ac 24 80 00 00 00     	leaq	0x80(%rsp), %rbp
14008ab83: 4c 8b 4c 24 30              	movq	0x30(%rsp), %r9
14008ab88: 66 44 89 04 01              	movw	%r8w, (%rcx,%rax)
14008ab8d: 48 8d 84 24 40 02 00 00     	leaq	0x240(%rsp), %rax
14008ab95: 49 89 e8                    	movq	%rbp, %r8
14008ab98: 48 8d 4c 24 60              	leaq	0x60(%rsp), %rcx
14008ab9d: 48 89 44 24 20              	movq	%rax, 0x20(%rsp)
14008aba2: 48 89 94 24 68 02 00 00     	movq	%rdx, 0x268(%rsp)
14008abaa: 48 8b 54 24 38              	movq	0x38(%rsp), %rdx
14008abaf: e8 3c eb ff ff              	callq	0x1400896f0 <.text+0x886f0>
14008abb4: 48 8b 8c 24 60 02 00 00     	movq	0x260(%rsp), %rcx
14008abbc: 48 39 f9                    	cmpq	%rdi, %rcx
14008abbf: 74 12                       	je	0x14008abd3 <.text+0x89bd3>
14008abc1: 48 8b 84 24 70 02 00 00     	movq	0x270(%rsp), %rax
14008abc9: 48 8d 54 00 02              	leaq	0x2(%rax,%rax), %rdx
14008abce: e8 fd f4 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008abd3: 48 8b 4c 24 30              	movq	0x30(%rsp), %rcx
14008abd8: e8 a3 ea ff ff              	callq	0x140089680 <.text+0x88680>
14008abdd: 84 c0                       	testb	%al, %al
14008abdf: 0f 85 77 05 00 00           	jne	0x14008b15c <.text+0x8a15c>
14008abe5: 48 8b 8c 24 40 02 00 00     	movq	0x240(%rsp), %rcx
14008abed: 48 39 f1                    	cmpq	%rsi, %rcx
14008abf0: 74 12                       	je	0x14008ac04 <.text+0x89c04>
14008abf2: 48 8b 84 24 50 02 00 00     	movq	0x250(%rsp), %rax
14008abfa: 48 8d 54 00 02              	leaq	0x2(%rax,%rax), %rdx
14008abff: e8 cc f4 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008ac04: 48 8b 8c 24 20 02 00 00     	movq	0x220(%rsp), %rcx
14008ac0c: 48 39 d9                    	cmpq	%rbx, %rcx
14008ac0f: 74 11                       	je	0x14008ac22 <.text+0x89c22>
14008ac11: 48 8b 84 24 30 02 00 00     	movq	0x230(%rsp), %rax
14008ac19: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008ac1d: e8 ae f4 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008ac22: 48 8b 8c 24 80 00 00 00     	movq	0x80(%rsp), %rcx
14008ac2a: 4c 39 f9                    	cmpq	%r15, %rcx
14008ac2d: 74 11                       	je	0x14008ac40 <.text+0x89c40>
14008ac2f: 48 8b 84 24 90 00 00 00     	movq	0x90(%rsp), %rax
14008ac37: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008ac3b: e8 90 f4 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008ac40: 48 8b 4c 24 60              	movq	0x60(%rsp), %rcx
14008ac45: 48 8d 44 24 70              	leaq	0x70(%rsp), %rax
14008ac4a: 48 39 c1                    	cmpq	%rax, %rcx
14008ac4d: 74 10                       	je	0x14008ac5f <.text+0x89c5f>
14008ac4f: 48 8b 44 24 70              	movq	0x70(%rsp), %rax
14008ac54: 48 8d 54 00 02              	leaq	0x2(%rax,%rax), %rdx
14008ac59: e8 72 f4 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008ac5e: 90                          	nop
14008ac5f: 48 81 c4 38 03 00 00        	addq	$0x338, %rsp            # imm = 0x338
14008ac66: 5b                          	popq	%rbx
14008ac67: 5e                          	popq	%rsi
14008ac68: 5f                          	popq	%rdi
14008ac69: 5d                          	popq	%rbp
14008ac6a: 41 5c                       	popq	%r12
14008ac6c: 41 5d                       	popq	%r13
14008ac6e: 41 5e                       	popq	%r14
14008ac70: 41 5f                       	popq	%r15
14008ac72: c3                          	retq
14008ac73: 0f 1f 44 00 00              	nopl	(%rax,%rax)
14008ac78: 41 b9 2f 00 00 00           	movl	$0x2f, %r9d
14008ac7e: b8 02 00 00 00              	movl	$0x2, %eax
14008ac83: 66 44 89 09                 	movw	%r9w, (%rcx)
14008ac87: e9 ec fe ff ff              	jmp	0x14008ab78 <.text+0x89b78>
14008ac8c: 0f 1f 40 00                 	nopl	(%rax)
14008ac90: 48 8b 54 24 50              	movq	0x50(%rsp), %rdx
14008ac95: 48 8b 4c 24 38              	movq	0x38(%rsp), %rcx
14008ac9a: 45 31 c0                    	xorl	%r8d, %r8d
14008ac9d: 48 8d ac 24 80 00 00 00     	leaq	0x80(%rsp), %rbp
14008aca5: e8 36 bb 03 00              	callq	0x1400c67e0 <.text+0xc57e0>
14008acaa: 48 89 84 24 60 02 00 00     	movq	%rax, 0x260(%rsp)
14008acb2: 48 89 c1                    	movq	%rax, %rcx
14008acb5: 48 8b 84 24 00 02 00 00     	movq	0x200(%rsp), %rax
14008acbd: 48 89 84 24 70 02 00 00     	movq	%rax, 0x270(%rsp)
14008acc5: 48 8d 15 4c 00 07 00        	leaq	0x7004c(%rip), %rdx     # 0x1400fad18
14008accc: 4d 89 e0                    	movq	%r12, %r8
14008accf: e8 14 ab 01 00              	callq	0x1400a57e8 <.text+0xa47e8>
14008acd4: 48 8b 94 24 00 02 00 00     	movq	0x200(%rsp), %rdx
14008acdc: 48 8b 8c 24 60 02 00 00     	movq	0x260(%rsp), %rcx
14008ace4: 48 8d 04 12                 	leaq	(%rdx,%rdx), %rax
14008ace8: e9 8b fe ff ff              	jmp	0x14008ab78 <.text+0x89b78>
14008aced: 0f 1f 00                    	nopl	(%rax)
14008acf0: 45 8d 6e 02                 	leal	0x2(%r14), %r13d
14008acf4: 4d 89 ee                    	movq	%r13, %r14
14008acf7: e9 71 f4 ff ff              	jmp	0x14008a16d <.text+0x8916d>
14008acfc: 0f 1f 40 00                 	nopl	(%rax)
14008ad00: 45 8d 6e 03                 	leal	0x3(%r14), %r13d
14008ad04: 4d 89 ee                    	movq	%r13, %r14
14008ad07: e9 61 f4 ff ff              	jmp	0x14008a16d <.text+0x8916d>
14008ad0c: 0f 1f 40 00                 	nopl	(%rax)
14008ad10: 45 89 f5                    	movl	%r14d, %r13d
14008ad13: e9 55 f4 ff ff              	jmp	0x14008a16d <.text+0x8916d>
14008ad18: 0f 1f 84 00 00 00 00 00     	nopl	(%rax,%rax)
14008ad20: 48 89 54 24 20              	movq	%rdx, 0x20(%rsp)
14008ad25: 45 31 c0                    	xorl	%r8d, %r8d
14008ad28: 31 d2                       	xorl	%edx, %edx
14008ad2a: 48 8d 8c 24 60 01 00 00     	leaq	0x160(%rsp), %rcx
14008ad32: 48 8d 9c 24 20 01 00 00     	leaq	0x120(%rsp), %rbx
14008ad3a: e8 91 4b 03 00              	callq	0x1400bf8d0 <.text+0xbe8d0>
14008ad3f: 48 8d bc 24 90 01 00 00     	leaq	0x190(%rsp), %rdi
14008ad47: 48 8d 50 10                 	leaq	0x10(%rax), %rdx
14008ad4b: 48 89 bc 24 80 01 00 00     	movq	%rdi, 0x180(%rsp)
14008ad53: 48 8b 08                    	movq	(%rax), %rcx
14008ad56: 48 89 7c 24 40              	movq	%rdi, 0x40(%rsp)
14008ad5b: 48 39 d1                    	cmpq	%rdx, %rcx
14008ad5e: 0f 85 31 f8 ff ff           	jne	0x14008a595 <.text+0x89595>
14008ad64: 48 8b 48 08                 	movq	0x8(%rax), %rcx
14008ad68: 4c 8d 41 01                 	leaq	0x1(%rcx), %r8
14008ad6c: 41 83 e0 3f                 	andl	$0x3f, %r8d
14008ad70: 0f 84 37 f8 ff ff           	je	0x14008a5ad <.text+0x895ad>
14008ad76: 31 c9                       	xorl	%ecx, %ecx
14008ad78: 41 89 c9                    	movl	%ecx, %r9d
14008ad7b: 48 8b 7c 24 40              	movq	0x40(%rsp), %rdi
14008ad80: 83 c1 01                    	addl	$0x1, %ecx
14008ad83: 46 0f b6 14 0a              	movzbl	(%rdx,%r9), %r10d
14008ad88: 46 88 14 0f                 	movb	%r10b, (%rdi,%r9)
14008ad8c: 44 39 c1                    	cmpl	%r8d, %ecx
14008ad8f: 72 e7                       	jb	0x14008ad78 <.text+0x89d78>
14008ad91: e9 13 f8 ff ff              	jmp	0x14008a5a9 <.text+0x895a9>
14008ad96: 66 2e 0f 1f 84 00 00 00 00 00       	nopw	%cs:(%rax,%rax)
14008ada0: 48 89 44 24 20              	movq	%rax, 0x20(%rsp)
14008ada5: 45 31 c0                    	xorl	%r8d, %r8d
14008ada8: 31 d2                       	xorl	%edx, %edx
14008adaa: 48 8d 8c 24 e0 00 00 00     	leaq	0xe0(%rsp), %rcx
14008adb2: e8 19 4b 03 00              	callq	0x1400bf8d0 <.text+0xbe8d0>
14008adb7: 4c 8d ac 24 10 01 00 00     	leaq	0x110(%rsp), %r13
14008adbf: 48 8d 50 10                 	leaq	0x10(%rax), %rdx
14008adc3: 4c 89 ac 24 00 01 00 00     	movq	%r13, 0x100(%rsp)
14008adcb: 48 8b 08                    	movq	(%rax), %rcx
14008adce: 48 39 d1                    	cmpq	%rdx, %rcx
14008add1: 0f 85 6f f6 ff ff           	jne	0x14008a446 <.text+0x89446>
14008add7: 48 8b 48 08                 	movq	0x8(%rax), %rcx
14008addb: 4c 8d 41 01                 	leaq	0x1(%rcx), %r8
14008addf: 41 83 e0 3f                 	andl	$0x3f, %r8d
14008ade3: 0f 84 75 f6 ff ff           	je	0x14008a45e <.text+0x8945e>
14008ade9: 31 c9                       	xorl	%ecx, %ecx
14008adeb: 41 89 c9                    	movl	%ecx, %r9d
14008adee: 83 c1 01                    	addl	$0x1, %ecx
14008adf1: 46 0f b6 14 0a              	movzbl	(%rdx,%r9), %r10d
14008adf6: 47 88 54 0d 00              	movb	%r10b, (%r13,%r9)
14008adfb: 44 39 c1                    	cmpl	%r8d, %ecx
14008adfe: 72 eb                       	jb	0x14008adeb <.text+0x89deb>
14008ae00: e9 55 f6 ff ff              	jmp	0x14008a45a <.text+0x8945a>
14008ae05: 0f 1f 00                    	nopl	(%rax)
14008ae08: 48 89 44 24 20              	movq	%rax, 0x20(%rsp)
14008ae0d: 45 31 c0                    	xorl	%r8d, %r8d
14008ae10: 48 8d 84 24 20 02 00 00     	leaq	0x220(%rsp), %rax
14008ae18: 31 d2                       	xorl	%edx, %edx
14008ae1a: 48 8d 8c 24 40 02 00 00     	leaq	0x240(%rsp), %rcx
14008ae22: 48 89 44 24 30              	movq	%rax, 0x30(%rsp)
14008ae27: e8 a4 4a 03 00              	callq	0x1400bf8d0 <.text+0xbe8d0>
14008ae2c: 48 8d bc 24 70 02 00 00     	leaq	0x270(%rsp), %rdi
14008ae34: 48 8d 50 10                 	leaq	0x10(%rax), %rdx
14008ae38: 48 89 bc 24 60 02 00 00     	movq	%rdi, 0x260(%rsp)
14008ae40: 48 8b 08                    	movq	(%rax), %rcx
14008ae43: 48 39 d1                    	cmpq	%rdx, %rcx
14008ae46: 0f 84 d2 02 00 00           	je	0x14008b11e <.text+0x8a11e>
14008ae4c: 48 89 8c 24 60 02 00 00     	movq	%rcx, 0x260(%rsp)
14008ae54: 48 8b 48 10                 	movq	0x10(%rax), %rcx
14008ae58: 48 89 8c 24 70 02 00 00     	movq	%rcx, 0x270(%rsp)
14008ae60: 48 8b 48 08                 	movq	0x8(%rax), %rcx
14008ae64: 48 89 8c 24 68 02 00 00     	movq	%rcx, 0x268(%rsp)
14008ae6c: 48 89 10                    	movq	%rdx, (%rax)
14008ae6f: c6 40 10 00                 	movb	$0x0, 0x10(%rax)
14008ae73: 48 c7 40 08 00 00 00 00     	movq	$0x0, 0x8(%rax)
14008ae7b: 48 8d 84 24 20 02 00 00     	leaq	0x220(%rsp), %rax
14008ae83: 48 89 44 24 30              	movq	%rax, 0x30(%rsp)
14008ae88: e9 df f9 ff ff              	jmp	0x14008a86c <.text+0x8986c>
14008ae8d: 0f 1f 00                    	nopl	(%rax)
14008ae90: 48 89 54 24 20              	movq	%rdx, 0x20(%rsp)
14008ae95: 45 31 c0                    	xorl	%r8d, %r8d
14008ae98: 31 d2                       	xorl	%edx, %edx
14008ae9a: 48 8d 8c 24 e0 01 00 00     	leaq	0x1e0(%rsp), %rcx
14008aea2: 48 8d 9c 24 a0 01 00 00     	leaq	0x1a0(%rsp), %rbx
14008aeaa: e8 21 4a 03 00              	callq	0x1400bf8d0 <.text+0xbe8d0>
14008aeaf: 48 8d bc 24 10 02 00 00     	leaq	0x210(%rsp), %rdi
14008aeb7: 48 8d 50 10                 	leaq	0x10(%rax), %rdx
14008aebb: 48 89 bc 24 00 02 00 00     	movq	%rdi, 0x200(%rsp)
14008aec3: 48 8b 08                    	movq	(%rax), %rcx
14008aec6: 48 89 7c 24 48              	movq	%rdi, 0x48(%rsp)
14008aecb: 48 39 d1                    	cmpq	%rdx, %rcx
14008aece: 0f 85 15 f8 ff ff           	jne	0x14008a6e9 <.text+0x896e9>
14008aed4: 48 8b 48 08                 	movq	0x8(%rax), %rcx
14008aed8: 4c 8d 41 01                 	leaq	0x1(%rcx), %r8
14008aedc: 41 83 e0 3f                 	andl	$0x3f, %r8d
14008aee0: 0f 84 1b f8 ff ff           	je	0x14008a701 <.text+0x89701>
14008aee6: 31 c9                       	xorl	%ecx, %ecx
14008aee8: 41 89 c9                    	movl	%ecx, %r9d
14008aeeb: 48 8b 7c 24 48              	movq	0x48(%rsp), %rdi
14008aef0: 83 c1 01                    	addl	$0x1, %ecx
14008aef3: 46 0f b6 14 0a              	movzbl	(%rdx,%r9), %r10d
14008aef8: 46 88 14 0f                 	movb	%r10b, (%rdi,%r9)
14008aefc: 44 39 c1                    	cmpl	%r8d, %ecx
14008aeff: 72 e7                       	jb	0x14008aee8 <.text+0x89ee8>
14008af01: e9 f7 f7 ff ff              	jmp	0x14008a6fd <.text+0x896fd>
14008af06: 66 2e 0f 1f 84 00 00 00 00 00       	nopw	%cs:(%rax,%rax)
14008af10: 48 8b 48 08                 	movq	0x8(%rax), %rcx
14008af14: 4c 8d 41 01                 	leaq	0x1(%rcx), %r8
14008af18: 41 83 e0 3f                 	andl	$0x3f, %r8d
14008af1c: 0f 84 c0 f9 ff ff           	je	0x14008a8e2 <.text+0x898e2>
14008af22: 31 c9                       	xorl	%ecx, %ecx
14008af24: 41 89 c9                    	movl	%ecx, %r9d
14008af27: 83 c1 01                    	addl	$0x1, %ecx
14008af2a: 46 0f b6 14 0a              	movzbl	(%rdx,%r9), %r10d
14008af2f: 47 88 14 0f                 	movb	%r10b, (%r15,%r9)
14008af33: 44 39 c1                    	cmpl	%r8d, %ecx
14008af36: 72 ec                       	jb	0x14008af24 <.text+0x89f24>
14008af38: e9 a1 f9 ff ff              	jmp	0x14008a8de <.text+0x898de>
14008af3d: 0f 1f 00                    	nopl	(%rax)
14008af40: 48 83 fa 0f                 	cmpq	$0xf, %rdx
14008af44: 0f 86 bb f8 ff ff           	jbe	0x14008a805 <.text+0x89805>
14008af4a: 48 39 f1                    	cmpq	%rsi, %rcx
14008af4d: 0f 85 8e f8 ff ff           	jne	0x14008a7e1 <.text+0x897e1>
14008af53: e9 ad f8 ff ff              	jmp	0x14008a805 <.text+0x89805>
14008af58: 0f 1f 84 00 00 00 00 00     	nopl	(%rax,%rax)
14008af60: 48 8b 48 08                 	movq	0x8(%rax), %rcx
14008af64: 4c 8d 41 01                 	leaq	0x1(%rcx), %r8
14008af68: 41 83 e0 3f                 	andl	$0x3f, %r8d
14008af6c: 0f 84 1c f8 ff ff           	je	0x14008a78e <.text+0x8978e>
14008af72: 31 c9                       	xorl	%ecx, %ecx
14008af74: 41 89 c9                    	movl	%ecx, %r9d
14008af77: 83 c1 01                    	addl	$0x1, %ecx
14008af7a: 46 0f b6 14 0a              	movzbl	(%rdx,%r9), %r10d
14008af7f: 46 88 14 0b                 	movb	%r10b, (%rbx,%r9)
14008af83: 44 39 c1                    	cmpl	%r8d, %ecx
14008af86: 72 ec                       	jb	0x14008af74 <.text+0x89f74>
14008af88: e9 fd f7 ff ff              	jmp	0x14008a78a <.text+0x8978a>
14008af8d: 0f 1f 00                    	nopl	(%rax)
14008af90: 48 83 f8 0f                 	cmpq	$0xf, %rax
14008af94: 0f 86 c3 f5 ff ff           	jbe	0x14008a55d <.text+0x8955d>
14008af9a: 4c 8d 94 24 70 01 00 00     	leaq	0x170(%rsp), %r10
14008afa2: 4c 39 d1                    	cmpq	%r10, %rcx
14008afa5: 0f 85 8e f5 ff ff           	jne	0x14008a539 <.text+0x89539>
14008afab: e9 ad f5 ff ff              	jmp	0x14008a55d <.text+0x8955d>
14008afb0: 48 8b 50 08                 	movq	0x8(%rax), %rdx
14008afb4: 4c 8d 42 01                 	leaq	0x1(%rdx), %r8
14008afb8: 41 83 e0 3f                 	andl	$0x3f, %r8d
14008afbc: 0f 84 24 f5 ff ff           	je	0x14008a4e6 <.text+0x894e6>
14008afc2: 31 d2                       	xorl	%edx, %edx
14008afc4: 41 89 d1                    	movl	%edx, %r9d
14008afc7: 83 c2 01                    	addl	$0x1, %edx
14008afca: 46 0f b6 14 09              	movzbl	(%rcx,%r9), %r10d
14008afcf: 46 88 54 0d 00              	movb	%r10b, (%rbp,%r9)
14008afd4: 44 39 c2                    	cmpl	%r8d, %edx
14008afd7: 72 eb                       	jb	0x14008afc4 <.text+0x89fc4>
14008afd9: e9 04 f5 ff ff              	jmp	0x14008a4e2 <.text+0x894e2>
14008afde: 66 90                       	nop
14008afe0: 48 83 fa 0f                 	cmpq	$0xf, %rdx
14008afe4: 0f 86 2c f4 ff ff           	jbe	0x14008a416 <.text+0x89416>
14008afea: 4c 8d 94 24 f0 00 00 00     	leaq	0xf0(%rsp), %r10
14008aff2: 4c 39 d1                    	cmpq	%r10, %rcx
14008aff5: 0f 85 f7 f3 ff ff           	jne	0x14008a3f2 <.text+0x893f2>
14008affb: e9 16 f4 ff ff              	jmp	0x14008a416 <.text+0x89416>
14008b000: 48 83 f8 0f                 	cmpq	$0xf, %rax
14008b004: 0f 86 a7 f6 ff ff           	jbe	0x14008a6b1 <.text+0x896b1>
14008b00a: 4c 8d 94 24 f0 01 00 00     	leaq	0x1f0(%rsp), %r10
14008b012: 4c 39 d1                    	cmpq	%r10, %rcx
14008b015: 0f 85 72 f6 ff ff           	jne	0x14008a68d <.text+0x8968d>
14008b01b: e9 91 f6 ff ff              	jmp	0x14008a6b1 <.text+0x896b1>
14008b020: 48 8b 50 08                 	movq	0x8(%rax), %rdx
14008b024: 4c 8d 42 01                 	leaq	0x1(%rdx), %r8
14008b028: 41 83 e0 3f                 	andl	$0x3f, %r8d
14008b02c: 0f 84 08 f6 ff ff           	je	0x14008a63a <.text+0x8963a>
14008b032: 31 d2                       	xorl	%edx, %edx
14008b034: 41 89 d1                    	movl	%edx, %r9d
14008b037: 83 c2 01                    	addl	$0x1, %edx
14008b03a: 46 0f b6 14 09              	movzbl	(%rcx,%r9), %r10d
14008b03f: 47 88 14 0c                 	movb	%r10b, (%r12,%r9)
14008b043: 44 39 c2                    	cmpl	%r8d, %edx
14008b046: 72 ec                       	jb	0x14008b034 <.text+0x8a034>
14008b048: e9 e9 f5 ff ff              	jmp	0x14008a636 <.text+0x89636>
14008b04d: 0f 1f 00                    	nopl	(%rax)
14008b050: 48 8b 48 08                 	movq	0x8(%rax), %rcx
14008b054: 4c 8d 41 01                 	leaq	0x1(%rcx), %r8
14008b058: 41 83 e0 3f                 	andl	$0x3f, %r8d
14008b05c: 0f 84 9f f6 ff ff           	je	0x14008a701 <.text+0x89701>
14008b062: 31 c9                       	xorl	%ecx, %ecx
14008b064: 41 89 c9                    	movl	%ecx, %r9d
14008b067: 48 8b 7c 24 48              	movq	0x48(%rsp), %rdi
14008b06c: 83 c1 01                    	addl	$0x1, %ecx
14008b06f: 46 0f b6 14 0a              	movzbl	(%rdx,%r9), %r10d
14008b074: 46 88 14 0f                 	movb	%r10b, (%rdi,%r9)
14008b078: 44 39 c1                    	cmpl	%r8d, %ecx
14008b07b: 72 e7                       	jb	0x14008b064 <.text+0x8a064>
14008b07d: e9 7b f6 ff ff              	jmp	0x14008a6fd <.text+0x896fd>
14008b082: 66 0f 1f 44 00 00           	nopw	(%rax,%rax)
14008b088: 48 8b 48 08                 	movq	0x8(%rax), %rcx
14008b08c: 4c 8d 41 01                 	leaq	0x1(%rcx), %r8
14008b090: 41 83 e0 3f                 	andl	$0x3f, %r8d
14008b094: 0f 84 13 f5 ff ff           	je	0x14008a5ad <.text+0x895ad>
14008b09a: 31 c9                       	xorl	%ecx, %ecx
14008b09c: 41 89 c9                    	movl	%ecx, %r9d
14008b09f: 48 8b 7c 24 40              	movq	0x40(%rsp), %rdi
14008b0a4: 83 c1 01                    	addl	$0x1, %ecx
14008b0a7: 46 0f b6 14 0a              	movzbl	(%rdx,%r9), %r10d
14008b0ac: 46 88 14 0f                 	movb	%r10b, (%rdi,%r9)
14008b0b0: 44 39 c1                    	cmpl	%r8d, %ecx
14008b0b3: 72 e7                       	jb	0x14008b09c <.text+0x8a09c>
14008b0b5: e9 ef f4 ff ff              	jmp	0x14008a5a9 <.text+0x895a9>
14008b0ba: 66 0f 1f 44 00 00           	nopw	(%rax,%rax)
14008b0c0: 48 8b 48 08                 	movq	0x8(%rax), %rcx
14008b0c4: 4c 8d 41 01                 	leaq	0x1(%rcx), %r8
14008b0c8: 41 83 e0 3f                 	andl	$0x3f, %r8d
14008b0cc: 0f 84 83 f7 ff ff           	je	0x14008a855 <.text+0x89855>
14008b0d2: 31 c9                       	xorl	%ecx, %ecx
14008b0d4: 41 89 c9                    	movl	%ecx, %r9d
14008b0d7: 83 c1 01                    	addl	$0x1, %ecx
14008b0da: 46 0f b6 14 0a              	movzbl	(%rdx,%r9), %r10d
14008b0df: 46 88 14 0f                 	movb	%r10b, (%rdi,%r9)
14008b0e3: 44 39 c1                    	cmpl	%r8d, %ecx
14008b0e6: 72 ec                       	jb	0x14008b0d4 <.text+0x8a0d4>
14008b0e8: e9 64 f7 ff ff              	jmp	0x14008a851 <.text+0x89851>
14008b0ed: 0f 1f 00                    	nopl	(%rax)
14008b0f0: 48 8b 48 08                 	movq	0x8(%rax), %rcx
14008b0f4: 4c 8d 41 01                 	leaq	0x1(%rcx), %r8
14008b0f8: 41 83 e0 3f                 	andl	$0x3f, %r8d
14008b0fc: 0f 84 5c f3 ff ff           	je	0x14008a45e <.text+0x8945e>
14008b102: 31 c9                       	xorl	%ecx, %ecx
14008b104: 41 89 c9                    	movl	%ecx, %r9d
14008b107: 83 c1 01                    	addl	$0x1, %ecx
14008b10a: 46 0f b6 14 0a              	movzbl	(%rdx,%r9), %r10d
14008b10f: 47 88 54 0d 00              	movb	%r10b, (%r13,%r9)
14008b114: 44 39 c1                    	cmpl	%r8d, %ecx
14008b117: 72 eb                       	jb	0x14008b104 <.text+0x8a104>
14008b119: e9 3c f3 ff ff              	jmp	0x14008a45a <.text+0x8945a>
14008b11e: 48 8b 48 08                 	movq	0x8(%rax), %rcx
14008b122: 4c 8d 41 01                 	leaq	0x1(%rcx), %r8
14008b126: 41 83 e0 3f                 	andl	$0x3f, %r8d
14008b12a: 0f 84 34 fd ff ff           	je	0x14008ae64 <.text+0x89e64>
14008b130: 31 c9                       	xorl	%ecx, %ecx
14008b132: 41 89 c9                    	movl	%ecx, %r9d
14008b135: 83 c1 01                    	addl	$0x1, %ecx
14008b138: 46 0f b6 14 0a              	movzbl	(%rdx,%r9), %r10d
14008b13d: 46 88 14 0f                 	movb	%r10b, (%rdi,%r9)
14008b141: 44 39 c1                    	cmpl	%r8d, %ecx
14008b144: 72 ec                       	jb	0x14008b132 <.text+0x8a132>
14008b146: e9 15 fd ff ff              	jmp	0x14008ae60 <.text+0x89e60>
14008b14b: 41 bd 01 00 00 00           	movl	$0x1, %r13d
14008b151: 41 be 01 00 00 00           	movl	$0x1, %r14d
14008b157: e9 11 f0 ff ff              	jmp	0x14008a16d <.text+0x8916d>
14008b15c: 31 c9                       	xorl	%ecx, %ecx
14008b15e: ff 15 44 ac 0e 00           	callq	*0xeac44(%rip)          # 0x140175da8
14008b164: 48 89 c3                    	movq	%rax, %rbx
14008b167: e9 83 d0 04 00              	jmp	0x1400d81ef <.text+0xd71ef>
14008b16c: 48 89 c3                    	movq	%rax, %rbx
14008b16f: e9 88 d0 04 00              	jmp	0x1400d81fc <.text+0xd71fc>
14008b174: e9 11 d1 04 00              	jmp	0x1400d828a <.text+0xd728a>
14008b179: e9 2a d1 04 00              	jmp	0x1400d82a8 <.text+0xd72a8>
14008b17e: e9 57 d1 04 00              	jmp	0x1400d82da <.text+0xd72da>
14008b183: e9 57 d1 04 00              	jmp	0x1400d82df <.text+0xd72df>
14008b188: 48 89 c3                    	movq	%rax, %rbx
14008b18b: e9 6c d0 04 00              	jmp	0x1400d81fc <.text+0xd71fc>
14008b190: 48 89 c3                    	movq	%rax, %rbx
14008b193: e9 56 d1 04 00              	jmp	0x1400d82ee <.text+0xd72ee>
14008b198: 48 89 c3                    	movq	%rax, %rbx
14008b19b: e9 44 d1 04 00              	jmp	0x1400d82e4 <.text+0xd72e4>
14008b1a0: e9 8f d0 04 00              	jmp	0x1400d8234 <.text+0xd7234>
14008b1a5: 48 89 c3                    	movq	%rax, %rbx
14008b1a8: e9 0e d0 04 00              	jmp	0x1400d81bb <.text+0xd71bb>
14008b1ad: 48 89 c3                    	movq	%rax, %rbx
14008b1b0: e9 13 d0 04 00              	jmp	0x1400d81c8 <.text+0xd71c8>
14008b1b5: 48 89 c3                    	movq	%rax, %rbx
14008b1b8: e9 18 d0 04 00              	jmp	0x1400d81d5 <.text+0xd71d5>
14008b1bd: 48 89 c3                    	movq	%rax, %rbx
14008b1c0: e9 1d d0 04 00              	jmp	0x1400d81e2 <.text+0xd71e2>
14008b1c5: e9 61 d1 04 00              	jmp	0x1400d832b <.text+0xd732b>
14008b1ca: e9 c7 cf 04 00              	jmp	0x1400d8196 <.text+0xd7196>
14008b1cf: e9 87 d1 04 00              	jmp	0x1400d835b <.text+0xd735b>
14008b1d4: 48 89 c3                    	movq	%rax, %rbx
14008b1d7: e9 d2 cf 04 00              	jmp	0x1400d81ae <.text+0xd71ae>
14008b1dc: 0f 1f 40 00                 	nopl	(%rax)
14008b1e0: 41 57                       	pushq	%r15
14008b1e2: 41 56                       	pushq	%r14
14008b1e4: 41 55                       	pushq	%r13
14008b1e6: 41 54                       	pushq	%r12
14008b1e8: 55                          	pushq	%rbp
14008b1e9: 57                          	pushq	%rdi
14008b1ea: 56                          	pushq	%rsi
14008b1eb: 53                          	pushq	%rbx
14008b1ec: 48 81 ec 38 01 00 00        	subq	$0x138, %rsp            # imm = 0x138
14008b1f3: 31 c0                       	xorl	%eax, %eax
14008b1f5: 45 31 c0                    	xorl	%r8d, %r8d
14008b1f8: 66 89 01                    	movw	%ax, (%rcx)
14008b1fb: 31 c0                       	xorl	%eax, %eax
14008b1fd: 48 89 cb                    	movq	%rcx, %rbx
14008b200: 4c 8d 71 18                 	leaq	0x18(%rcx), %r14
14008b204: 66 89 41 18                 	movw	%ax, 0x18(%rcx)
14008b208: 31 c0                       	xorl	%eax, %eax
14008b20a: 4c 8d 69 50                 	leaq	0x50(%rcx), %r13
14008b20e: 4c 8d 61 70                 	leaq	0x70(%rcx), %r12
14008b212: 66 89 41 50                 	movw	%ax, 0x50(%rcx)
14008b216: 48 8d 83 90 00 00 00        	leaq	0x90(%rbx), %rax
14008b21d: 48 8d bc 24 80 00 00 00     	leaq	0x80(%rsp), %rdi
14008b225: 4c 89 71 08                 	movq	%r14, 0x8(%rcx)
14008b229: 48 c7 41 10 00 00 00 00     	movq	$0x0, 0x10(%rcx)
14008b231: 48 c7 41 28 00 00 00 00     	movq	$0x0, 0x28(%rcx)
14008b239: 48 c7 41 30 00 00 00 00     	movq	$0x0, 0x30(%rcx)
14008b241: 4c 89 69 40                 	movq	%r13, 0x40(%rcx)
14008b245: 48 c7 41 48 00 00 00 00     	movq	$0x0, 0x48(%rcx)
14008b24d: 4c 89 61 60                 	movq	%r12, 0x60(%rcx)
14008b251: c7 41 38 00 00 00 00        	movl	$0x0, 0x38(%rcx)
14008b258: c6 41 3c 00                 	movb	$0x0, 0x3c(%rcx)
14008b25c: 31 c9                       	xorl	%ecx, %ecx
14008b25e: 48 89 83 80 00 00 00        	movq	%rax, 0x80(%rbx)
14008b265: 48 8d 83 b0 00 00 00        	leaq	0xb0(%rbx), %rax
14008b26c: c6 83 b0 00 00 00 00        	movb	$0x0, 0xb0(%rbx)
14008b273: 66 89 4b 70                 	movw	%cx, 0x70(%rbx)
14008b277: 48 8d 4c 24 70              	leaq	0x70(%rsp), %rcx
14008b27c: 48 c7 43 68 00 00 00 00     	movq	$0x0, 0x68(%rbx)
14008b284: 66 44 89 83 90 00 00 00     	movw	%r8w, 0x90(%rbx)
14008b28c: 48 c7 83 88 00 00 00 00 00 00 00    	movq	$0x0, 0x88(%rbx)
14008b297: 48 89 83 a0 00 00 00        	movq	%rax, 0xa0(%rbx)
14008b29e: 48 8b 02                    	movq	(%rdx), %rax
14008b2a1: 48 c7 83 a8 00 00 00 00 00 00 00    	movq	$0x0, 0xa8(%rbx)
14008b2ac: 4c 8b 42 08                 	movq	0x8(%rdx), %r8
14008b2b0: 48 89 c2                    	movq	%rax, %rdx
14008b2b3: 48 89 7c 24 70              	movq	%rdi, 0x70(%rsp)
14008b2b8: e8 53 8e 03 00              	callq	0x1400c4110 <.text+0xc3110>
14008b2bd: 48 8b 44 24 70              	movq	0x70(%rsp), %rax
14008b2c2: 48 8b 54 24 78              	movq	0x78(%rsp), %rdx
14008b2c7: 4c 8d 0c 50                 	leaq	(%rax,%rdx,2), %r9
14008b2cb: 4c 39 c8                    	cmpq	%r9, %rax
14008b2ce: 74 6b                       	je	0x14008b33b <.text+0x8a33b>
14008b2d0: 0f 1f 44 00 00              	nopl	(%rax,%rax)
14008b2d5: 66 66 2e 0f 1f 84 00 00 00 00 00    	nopw	%cs:(%rax,%rax)
14008b2e0: 0f b7 08                    	movzwl	(%rax), %ecx
14008b2e3: 44 8d 41 9f                 	leal	-0x61(%rcx), %r8d
14008b2e7: 66 41 83 f8 19              	cmpw	$0x19, %r8w
14008b2ec: 77 06                       	ja	0x14008b2f4 <.text+0x8a2f4>
14008b2ee: 83 e9 20                    	subl	$0x20, %ecx
14008b2f1: 66 89 08                    	movw	%cx, (%rax)
14008b2f4: 48 83 c0 02                 	addq	$0x2, %rax
14008b2f8: 4c 39 c8                    	cmpq	%r9, %rax
14008b2fb: 75 e3                       	jne	0x14008b2e0 <.text+0x8a2e0>
14008b2fd: 48 85 d2                    	testq	%rdx, %rdx
14008b300: 74 3e                       	je	0x14008b340 <.text+0x8a340>
14008b302: 0f 1f 00                    	nopl	(%rax)
14008b305: 66 66 2e 0f 1f 84 00 00 00 00 00    	nopw	%cs:(%rax,%rax)
14008b310: 48 8b 44 24 70              	movq	0x70(%rsp), %rax
14008b315: 0f b7 08                    	movzwl	(%rax), %ecx
14008b318: 66 83 f9 20                 	cmpw	$0x20, %cx
14008b31c: 74 06                       	je	0x14008b324 <.text+0x8a324>
14008b31e: 66 83 f9 09                 	cmpw	$0x9, %cx
14008b322: 75 7c                       	jne	0x14008b3a0 <.text+0x8a3a0>
14008b324: 41 b8 01 00 00 00           	movl	$0x1, %r8d
14008b32a: 31 d2                       	xorl	%edx, %edx
14008b32c: 48 8d 4c 24 70              	leaq	0x70(%rsp), %rcx
14008b331: e8 2a b2 03 00              	callq	0x1400c6560 <.text+0xc5560>
14008b336: 48 8b 54 24 78              	movq	0x78(%rsp), %rdx
14008b33b: 48 85 d2                    	testq	%rdx, %rdx
14008b33e: 75 d0                       	jne	0x14008b310 <.text+0x8a310>
14008b340: 48 8d 0d 1f fa 06 00        	leaq	0x6fa1f(%rip), %rcx     # 0x1400fad66
14008b347: 48 8d 73 08                 	leaq	0x8(%rbx), %rsi
14008b34b: e8 30 a3 01 00              	callq	0x1400a5680 <.text+0xa4680>
14008b350: 4c 8b 43 10                 	movq	0x10(%rbx), %r8
14008b354: 31 d2                       	xorl	%edx, %edx
14008b356: 48 89 f1                    	movq	%rsi, %rcx
14008b359: 48 89 44 24 20              	movq	%rax, 0x20(%rsp)
14008b35e: 4c 8d 0d 01 fa 06 00        	leaq	0x6fa01(%rip), %r9      # 0x1400fad66
14008b365: e8 b6 89 03 00              	callq	0x1400c3d20 <.text+0xc2d20>
14008b36a: 48 8b 4c 24 70              	movq	0x70(%rsp), %rcx
14008b36f: 48 39 f9                    	cmpq	%rdi, %rcx
14008b372: 74 12                       	je	0x14008b386 <.text+0x8a386>
14008b374: 48 8b 84 24 80 00 00 00     	movq	0x80(%rsp), %rax
14008b37c: 48 8d 54 00 02              	leaq	0x2(%rax,%rax), %rdx
14008b381: e8 4a ed 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008b386: 48 89 d8                    	movq	%rbx, %rax
14008b389: 48 81 c4 38 01 00 00        	addq	$0x138, %rsp            # imm = 0x138
14008b390: 5b                          	popq	%rbx
14008b391: 5e                          	popq	%rsi
14008b392: 5f                          	popq	%rdi
14008b393: 5d                          	popq	%rbp
14008b394: 41 5c                       	popq	%r12
14008b396: 41 5d                       	popq	%r13
14008b398: 41 5e                       	popq	%r14
14008b39a: 41 5f                       	popq	%r15
14008b39c: c3                          	retq
14008b39d: 0f 1f 00                    	nopl	(%rax)
14008b3a0: 48 89 d1                    	movq	%rdx, %rcx
14008b3a3: eb 08                       	jmp	0x14008b3ad <.text+0x8a3ad>
14008b3a5: 0f 1f 00                    	nopl	(%rax)
14008b3a8: 48 8b 44 24 70              	movq	0x70(%rsp), %rax
14008b3ad: 48 8d 51 ff                 	leaq	-0x1(%rcx), %rdx
14008b3b1: 44 0f b7 04 50              	movzwl	(%rax,%rdx,2), %r8d
14008b3b6: 66 41 83 f8 20              	cmpw	$0x20, %r8w
14008b3bb: 74 07                       	je	0x14008b3c4 <.text+0x8a3c4>
14008b3bd: 66 41 83 f8 09              	cmpw	$0x9, %r8w
14008b3c2: 75 24                       	jne	0x14008b3e8 <.text+0x8a3e8>
14008b3c4: 48 8d 4c 24 70              	leaq	0x70(%rsp), %rcx
14008b3c9: 41 b8 01 00 00 00           	movl	$0x1, %r8d
14008b3cf: e8 8c b1 03 00              	callq	0x1400c6560 <.text+0xc5560>
14008b3d4: 48 8b 4c 24 78              	movq	0x78(%rsp), %rcx
14008b3d9: 48 85 c9                    	testq	%rcx, %rcx
14008b3dc: 75 ca                       	jne	0x14008b3a8 <.text+0x8a3a8>
14008b3de: e9 5d ff ff ff              	jmp	0x14008b340 <.text+0x8a340>
14008b3e3: 0f 1f 44 00 00              	nopl	(%rax,%rax)
14008b3e8: 48 83 f9 10                 	cmpq	$0x10, %rcx
14008b3ec: 0f 85 4e ff ff ff           	jne	0x14008b340 <.text+0x8a340>
14008b3f2: 48 8d 48 20                 	leaq	0x20(%rax), %rcx
14008b3f6: 66 2e 0f 1f 84 00 00 00 00 00       	nopw	%cs:(%rax,%rax)
14008b400: 0f b7 30                    	movzwl	(%rax), %esi
14008b403: 8d 56 bf                    	leal	-0x41(%rsi), %edx
14008b406: 66 83 fa 19                 	cmpw	$0x19, %dx
14008b40a: 0f 87 b8 03 00 00           	ja	0x14008b7c8 <.text+0x8a7c8>
14008b410: 48 83 c0 02                 	addq	$0x2, %rax
14008b414: 48 39 c1                    	cmpq	%rax, %rcx
14008b417: 75 e7                       	jne	0x14008b400 <.text+0x8a400>
14008b419: 48 8d 54 24 70              	leaq	0x70(%rsp), %rdx
14008b41e: 48 8d 8c 24 f0 00 00 00     	leaq	0xf0(%rsp), %rcx
14008b426: e8 25 da ff ff              	callq	0x140088e50 <.text+0x87e50>
14008b42b: 4c 8b 84 24 f8 00 00 00     	movq	0xf8(%rsp), %r8
14008b433: 48 8d ac 24 e0 00 00 00     	leaq	0xe0(%rsp), %rbp
14008b43b: 48 b8 7b 22 63 61 72 64 5f 6b       	movabsq	$0x6b5f64726163227b, %rax # imm = 0x6B5F64726163227B
14008b445: c6 84 24 ed 00 00 00 00     	movb	$0x0, 0xed(%rsp)
14008b44d: 48 89 84 24 e0 00 00 00     	movq	%rax, 0xe0(%rsp)
14008b455: 48 8b 94 24 f0 00 00 00     	movq	0xf0(%rsp), %rdx
14008b45d: 48 b8 64 5f 6b 65 79 22 3a 22       	movabsq	$0x223a2279656b5f64, %rax # imm = 0x223A2279656B5F64
14008b467: 48 89 84 24 e5 00 00 00     	movq	%rax, 0xe5(%rsp)
14008b46f: 49 8d 40 0d                 	leaq	0xd(%r8), %rax
14008b473: 48 89 ac 24 d0 00 00 00     	movq	%rbp, 0xd0(%rsp)
14008b47b: 48 c7 84 24 d8 00 00 00 0d 00 00 00 	movq	$0xd, 0xd8(%rsp)
14008b487: 48 83 f8 0f                 	cmpq	$0xf, %rax
14008b48b: 76 2e                       	jbe	0x14008b4bb <.text+0x8a4bb>
14008b48d: 48 8d 8c 24 00 01 00 00     	leaq	0x100(%rsp), %rcx
14008b495: 48 39 ca                    	cmpq	%rcx, %rdx
14008b498: 74 21                       	je	0x14008b4bb <.text+0x8a4bb>
14008b49a: 48 39 84 24 00 01 00 00     	cmpq	%rax, 0x100(%rsp)
14008b4a2: 0f 83 7e 03 00 00           	jae	0x14008b826 <.text+0x8a826>
14008b4a8: 48 b8 f1 ff ff ff ff ff ff 7f       	movabsq	$0x7ffffffffffffff1, %rax # imm = 0x7FFFFFFFFFFFFFF1
14008b4b2: 4c 39 c0                    	cmpq	%r8, %rax
14008b4b5: 0f 82 09 cf 04 00           	jb	0x1400d83c4 <.text+0xd73c4>
14008b4bb: 48 8d 84 24 d0 00 00 00     	leaq	0xd0(%rsp), %rax
14008b4c3: 48 89 c1                    	movq	%rax, %rcx
14008b4c6: 48 89 44 24 40              	movq	%rax, 0x40(%rsp)
14008b4cb: e8 00 74 03 00              	callq	0x1400c28d0 <.text+0xc18d0>
14008b4d0: 48 8d b4 24 20 01 00 00     	leaq	0x120(%rsp), %rsi
14008b4d8: 48 8d 50 10                 	leaq	0x10(%rax), %rdx
14008b4dc: 48 89 b4 24 10 01 00 00     	movq	%rsi, 0x110(%rsp)
14008b4e4: 48 8b 08                    	movq	(%rax), %rcx
14008b4e7: 48 39 d1                    	cmpq	%rdx, %rcx
14008b4ea: 0f 84 93 05 00 00           	je	0x14008ba83 <.text+0x8aa83>
14008b4f0: 48 89 8c 24 10 01 00 00     	movq	%rcx, 0x110(%rsp)
14008b4f8: 48 8b 48 10                 	movq	0x10(%rax), %rcx
14008b4fc: 48 89 8c 24 20 01 00 00     	movq	%rcx, 0x120(%rsp)
14008b504: 48 8b 48 08                 	movq	0x8(%rax), %rcx
14008b508: 48 89 8c 24 18 01 00 00     	movq	%rcx, 0x118(%rsp)
14008b510: 48 89 10                    	movq	%rdx, (%rax)
14008b513: c6 40 10 00                 	movb	$0x0, 0x10(%rax)
14008b517: 48 c7 40 08 00 00 00 00     	movq	$0x0, 0x8(%rax)
14008b51f: 48 b8 03 00 00 00 00 00 00 80       	movabsq	$-0x7ffffffffffffffd, %rax # imm = 0x8000000000000003
14008b529: 48 03 84 24 18 01 00 00     	addq	0x118(%rsp), %rax
14008b531: 48 83 f8 01                 	cmpq	$0x1, %rax
14008b535: 0f 86 07 cf 04 00           	jbe	0x1400d8442 <.text+0xd7442>
14008b53b: 48 8d 84 24 10 01 00 00     	leaq	0x110(%rsp), %rax
14008b543: 41 b8 02 00 00 00           	movl	$0x2, %r8d
14008b549: 48 8d 15 c5 f7 06 00        	leaq	0x6f7c5(%rip), %rdx     # 0x1400fad15
14008b550: 48 89 c1                    	movq	%rax, %rcx
14008b553: 48 89 44 24 38              	movq	%rax, 0x38(%rsp)
14008b558: e8 73 73 03 00              	callq	0x1400c28d0 <.text+0xc18d0>
14008b55d: 48 8d 94 24 a0 00 00 00     	leaq	0xa0(%rsp), %rdx
14008b565: 48 89 94 24 90 00 00 00     	movq	%rdx, 0x90(%rsp)
14008b56d: 48 8b 08                    	movq	(%rax), %rcx
14008b570: 48 8d 50 10                 	leaq	0x10(%rax), %rdx
14008b574: 48 39 d1                    	cmpq	%rdx, %rcx
14008b577: 0f 84 d5 04 00 00           	je	0x14008ba52 <.text+0x8aa52>
14008b57d: 48 89 8c 24 90 00 00 00     	movq	%rcx, 0x90(%rsp)
14008b585: 48 8b 48 10                 	movq	0x10(%rax), %rcx
14008b589: 48 89 8c 24 a0 00 00 00     	movq	%rcx, 0xa0(%rsp)
14008b591: 48 8b 48 08                 	movq	0x8(%rax), %rcx
14008b595: 48 89 8c 24 98 00 00 00     	movq	%rcx, 0x98(%rsp)
14008b59d: 48 89 10                    	movq	%rdx, (%rax)
14008b5a0: c6 40 10 00                 	movb	$0x0, 0x10(%rax)
14008b5a4: 48 8b 8c 24 10 01 00 00     	movq	0x110(%rsp), %rcx
14008b5ac: 48 c7 40 08 00 00 00 00     	movq	$0x0, 0x8(%rax)
14008b5b4: 48 39 f1                    	cmpq	%rsi, %rcx
14008b5b7: 74 11                       	je	0x14008b5ca <.text+0x8a5ca>
14008b5b9: 48 8b 84 24 20 01 00 00     	movq	0x120(%rsp), %rax
14008b5c1: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008b5c5: e8 06 eb 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008b5ca: 48 8b 8c 24 d0 00 00 00     	movq	0xd0(%rsp), %rcx
14008b5d2: 48 39 e9                    	cmpq	%rbp, %rcx
14008b5d5: 74 11                       	je	0x14008b5e8 <.text+0x8a5e8>
14008b5d7: 48 8b 84 24 e0 00 00 00     	movq	0xe0(%rsp), %rax
14008b5df: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008b5e3: e8 e8 ea 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008b5e8: 48 8b 8c 24 f0 00 00 00     	movq	0xf0(%rsp), %rcx
14008b5f0: 4c 8d bc 24 00 01 00 00     	leaq	0x100(%rsp), %r15
14008b5f8: 4c 39 f9                    	cmpq	%r15, %rcx
14008b5fb: 74 11                       	je	0x14008b60e <.text+0x8a60e>
14008b5fd: 48 8b 84 24 00 01 00 00     	movq	0x100(%rsp), %rax
14008b605: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008b609: e8 c2 ea 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008b60e: 48 8d 84 24 c0 00 00 00     	leaq	0xc0(%rsp), %rax
14008b616: 48 8d 0d 5f f7 06 00        	leaq	0x6f75f(%rip), %rcx     # 0x1400fad7c
14008b61d: c6 84 24 c0 00 00 00 00     	movb	$0x0, 0xc0(%rsp)
14008b625: 48 89 84 24 b0 00 00 00     	movq	%rax, 0xb0(%rsp)
14008b62d: 31 c0                       	xorl	%eax, %eax
14008b62f: 48 c7 84 24 b8 00 00 00 00 00 00 00 	movq	$0x0, 0xb8(%rsp)
14008b63b: 48 89 ac 24 d0 00 00 00     	movq	%rbp, 0xd0(%rsp)
14008b643: 66 89 84 24 e0 00 00 00     	movw	%ax, 0xe0(%rsp)
14008b64b: 48 c7 84 24 d8 00 00 00 00 00 00 00 	movq	$0x0, 0xd8(%rsp)
14008b657: 48 89 b4 24 10 01 00 00     	movq	%rsi, 0x110(%rsp)
14008b65f: e8 1c a0 01 00              	callq	0x1400a5680 <.text+0xa4680>
14008b664: 4c 8d 0c 00                 	leaq	(%rax,%rax), %r9
14008b668: 4c 89 ca                    	movq	%r9, %rdx
14008b66b: 48 d1 fa                    	sarq	%rdx
14008b66e: 48 89 94 24 f0 00 00 00     	movq	%rdx, 0xf0(%rsp)
14008b676: 49 83 f9 0e                 	cmpq	$0xe, %r9
14008b67a: 0f 87 40 03 00 00           	ja	0x14008b9c0 <.text+0x8a9c0>
14008b680: 48 8b 8c 24 10 01 00 00     	movq	0x110(%rsp), %rcx
14008b688: 48 83 fa 01                 	cmpq	$0x1, %rdx
14008b68c: 0f 84 ae 03 00 00           	je	0x14008ba40 <.text+0x8aa40>
14008b692: 31 c0                       	xorl	%eax, %eax
14008b694: 48 85 d2                    	testq	%rdx, %rdx
14008b697: 0f 85 77 03 00 00           	jne	0x14008ba14 <.text+0x8aa14>
14008b69d: 45 31 db                    	xorl	%r11d, %r11d
14008b6a0: 66 44 89 1c 01              	movw	%r11w, (%rcx,%rax)
14008b6a5: 48 8d 8c 24 f0 00 00 00     	leaq	0xf0(%rsp), %rcx
14008b6ad: 48 89 94 24 18 01 00 00     	movq	%rdx, 0x118(%rsp)
14008b6b5: e8 56 e6 ff ff              	callq	0x140089d10 <.text+0x88d10>
14008b6ba: 48 8b 54 24 40              	movq	0x40(%rsp), %rdx
14008b6bf: 48 8d 84 24 90 00 00 00     	leaq	0x90(%rsp), %rax
14008b6c7: 4c 8d 8c 24 b0 00 00 00     	leaq	0xb0(%rsp), %r9
14008b6cf: 49 89 c0                    	movq	%rax, %r8
14008b6d2: 48 8d 8c 24 f0 00 00 00     	leaq	0xf0(%rsp), %rcx
14008b6da: 4c 89 4c 24 48              	movq	%r9, 0x48(%rsp)
14008b6df: 48 89 54 24 20              	movq	%rdx, 0x20(%rsp)
14008b6e4: 48 8b 54 24 38              	movq	0x38(%rsp), %rdx
14008b6e9: 48 89 44 24 58              	movq	%rax, 0x58(%rsp)
14008b6ee: e8 fd df ff ff              	callq	0x1400896f0 <.text+0x886f0>
14008b6f3: 48 8b 8c 24 f0 00 00 00     	movq	0xf0(%rsp), %rcx
14008b6fb: 88 44 24 50                 	movb	%al, 0x50(%rsp)
14008b6ff: 4c 39 f9                    	cmpq	%r15, %rcx
14008b702: 74 12                       	je	0x14008b716 <.text+0x8a716>
14008b704: 48 8b 84 24 00 01 00 00     	movq	0x100(%rsp), %rax
14008b70c: 48 8d 54 00 02              	leaq	0x2(%rax,%rax), %rdx
14008b711: e8 ba e9 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008b716: 48 8b 8c 24 10 01 00 00     	movq	0x110(%rsp), %rcx
14008b71e: 48 39 f1                    	cmpq	%rsi, %rcx
14008b721: 74 12                       	je	0x14008b735 <.text+0x8a735>
14008b723: 48 8b 84 24 20 01 00 00     	movq	0x120(%rsp), %rax
14008b72b: 48 8d 54 00 02              	leaq	0x2(%rax,%rax), %rdx
14008b730: e8 9b e9 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008b735: 80 7c 24 50 00              	cmpb	$0x0, 0x50(%rsp)
14008b73a: 0f 85 b8 00 00 00           	jne	0x14008b7f8 <.text+0x8a7f8>
14008b740: 48 8b 54 24 40              	movq	0x40(%rsp), %rdx
14008b745: 48 8d 4b 08                 	leaq	0x8(%rbx), %rcx
14008b749: e8 72 af 03 00              	callq	0x1400c66c0 <.text+0xc56c0>
14008b74e: 48 8b 8c 24 d0 00 00 00     	movq	0xd0(%rsp), %rcx
14008b756: 48 39 e9                    	cmpq	%rbp, %rcx
14008b759: 74 12                       	je	0x14008b76d <.text+0x8a76d>
14008b75b: 48 8b 84 24 e0 00 00 00     	movq	0xe0(%rsp), %rax
14008b763: 48 8d 54 00 02              	leaq	0x2(%rax,%rax), %rdx
14008b768: e8 63 e9 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008b76d: 48 8b 8c 24 b0 00 00 00     	movq	0xb0(%rsp), %rcx
14008b775: 48 8d 84 24 c0 00 00 00     	leaq	0xc0(%rsp), %rax
14008b77d: 48 39 c1                    	cmpq	%rax, %rcx
14008b780: 74 11                       	je	0x14008b793 <.text+0x8a793>
14008b782: 48 8b 84 24 c0 00 00 00     	movq	0xc0(%rsp), %rax
14008b78a: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008b78e: e8 3d e9 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008b793: 48 8b 8c 24 90 00 00 00     	movq	0x90(%rsp), %rcx
14008b79b: 48 8d 84 24 a0 00 00 00     	leaq	0xa0(%rsp), %rax
14008b7a3: 48 39 c1                    	cmpq	%rax, %rcx
14008b7a6: 0f 84 be fb ff ff           	je	0x14008b36a <.text+0x8a36a>
14008b7ac: 48 8b 84 24 a0 00 00 00     	movq	0xa0(%rsp), %rax
14008b7b4: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008b7b8: e8 13 e9 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008b7bd: e9 a8 fb ff ff              	jmp	0x14008b36a <.text+0x8a36a>
14008b7c2: 66 0f 1f 44 00 00           	nopw	(%rax,%rax)
14008b7c8: 48 8d 0d 97 f5 06 00        	leaq	0x6f597(%rip), %rcx     # 0x1400fad66
14008b7cf: e8 ac 9e 01 00              	callq	0x1400a5680 <.text+0xa4680>
14008b7d4: 48 8d 4b 08                 	leaq	0x8(%rbx), %rcx
14008b7d8: 4c 8d 0d 87 f5 06 00        	leaq	0x6f587(%rip), %r9      # 0x1400fad66
14008b7df: 31 d2                       	xorl	%edx, %edx
14008b7e1: 48 89 44 24 20              	movq	%rax, 0x20(%rsp)
14008b7e6: 4c 8b 43 10                 	movq	0x10(%rbx), %r8
14008b7ea: e8 31 85 03 00              	callq	0x1400c3d20 <.text+0xc2d20>
14008b7ef: e9 76 fb ff ff              	jmp	0x14008b36a <.text+0x8a36a>
14008b7f4: 0f 1f 40 00                 	nopl	(%rax)
14008b7f8: 48 8b 4c 24 48              	movq	0x48(%rsp), %rcx
14008b7fd: e8 7e de ff ff              	callq	0x140089680 <.text+0x88680>
14008b802: 84 c0                       	testb	%al, %al
14008b804: 0f 84 a8 00 00 00           	je	0x14008b8b2 <.text+0x8a8b2>
14008b80a: 48 8b 43 08                 	movq	0x8(%rbx), %rax
14008b80e: 45 31 d2                    	xorl	%r10d, %r10d
14008b811: c6 43 01 01                 	movb	$0x1, 0x1(%rbx)
14008b815: 66 44 89 10                 	movw	%r10w, (%rax)
14008b819: 48 c7 43 10 00 00 00 00     	movq	$0x0, 0x10(%rbx)
14008b821: e9 28 ff ff ff              	jmp	0x14008b74e <.text+0x8a74e>
14008b826: 48 8d 84 24 d0 00 00 00     	leaq	0xd0(%rsp), %rax
14008b82e: 49 89 e9                    	movq	%rbp, %r9
14008b831: 45 31 c0                    	xorl	%r8d, %r8d
14008b834: 31 d2                       	xorl	%edx, %edx
14008b836: 48 8d 8c 24 f0 00 00 00     	leaq	0xf0(%rsp), %rcx
14008b83e: 48 89 44 24 40              	movq	%rax, 0x40(%rsp)
14008b843: 48 c7 44 24 20 0d 00 00 00  	movq	$0xd, 0x20(%rsp)
14008b84c: e8 7f 40 03 00              	callq	0x1400bf8d0 <.text+0xbe8d0>
14008b851: 48 8d b4 24 20 01 00 00     	leaq	0x120(%rsp), %rsi
14008b859: 48 8d 50 10                 	leaq	0x10(%rax), %rdx
14008b85d: 48 89 b4 24 10 01 00 00     	movq	%rsi, 0x110(%rsp)
14008b865: 48 8b 08                    	movq	(%rax), %rcx
14008b868: 48 39 d1                    	cmpq	%rdx, %rcx
14008b86b: 0f 84 ad 05 00 00           	je	0x14008be1e <.text+0x8ae1e>
14008b871: 48 89 8c 24 10 01 00 00     	movq	%rcx, 0x110(%rsp)
14008b879: 48 8b 48 10                 	movq	0x10(%rax), %rcx
14008b87d: 48 89 8c 24 20 01 00 00     	movq	%rcx, 0x120(%rsp)
14008b885: 48 8b 48 08                 	movq	0x8(%rax), %rcx
14008b889: 48 89 8c 24 18 01 00 00     	movq	%rcx, 0x118(%rsp)
14008b891: 48 89 10                    	movq	%rdx, (%rax)
14008b894: c6 40 10 00                 	movb	$0x0, 0x10(%rax)
14008b898: 48 c7 40 08 00 00 00 00     	movq	$0x0, 0x8(%rax)
14008b8a0: 48 8d 84 24 d0 00 00 00     	leaq	0xd0(%rsp), %rax
14008b8a8: 48 89 44 24 40              	movq	%rax, 0x40(%rsp)
14008b8ad: e9 6d fc ff ff              	jmp	0x14008b51f <.text+0x8a51f>
14008b8b2: 48 8b 54 24 48              	movq	0x48(%rsp), %rdx
14008b8b7: 48 8b 4c 24 38              	movq	0x38(%rsp), %rcx
14008b8bc: 4c 8d 05 d7 f4 06 00        	leaq	0x6f4d7(%rip), %r8      # 0x1400fad9a
14008b8c3: e8 78 da ff ff              	callq	0x140089340 <.text+0x88340>
14008b8c8: 48 8b 8c 24 10 01 00 00     	movq	0x110(%rsp), %rcx
14008b8d0: 48 83 bc 24 18 01 00 00 04  	cmpq	$0x4, 0x118(%rsp)
14008b8d9: 0f 84 2e 02 00 00           	je	0x14008bb0d <.text+0x8ab0d>
14008b8df: 48 39 f1                    	cmpq	%rsi, %rcx
14008b8e2: 74 11                       	je	0x14008b8f5 <.text+0x8a8f5>
14008b8e4: 48 8b 84 24 20 01 00 00     	movq	0x120(%rsp), %rax
14008b8ec: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008b8f0: e8 db e7 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008b8f5: 48 8b 54 24 48              	movq	0x48(%rsp), %rdx
14008b8fa: 4c 8d 05 9c f4 06 00        	leaq	0x6f49c(%rip), %r8      # 0x1400fad9d
14008b901: 48 8d 8c 24 f0 00 00 00     	leaq	0xf0(%rsp), %rcx
14008b909: e8 32 da ff ff              	callq	0x140089340 <.text+0x88340>
14008b90e: 48 83 bc 24 f8 00 00 00 00  	cmpq	$0x0, 0xf8(%rsp)
14008b917: 0f 84 93 01 00 00           	je	0x14008bab0 <.text+0x8aab0>
14008b91d: 48 8b 4c 24 38              	movq	0x38(%rsp), %rcx
14008b922: 48 8d 94 24 f0 00 00 00     	leaq	0xf0(%rsp), %rdx
14008b92a: e8 a1 d7 ff ff              	callq	0x1400890d0 <.text+0x880d0>
14008b92f: 48 8b 84 24 10 01 00 00     	movq	0x110(%rsp), %rax
14008b937: 4c 8b 84 24 18 01 00 00     	movq	0x118(%rsp), %r8
14008b93f: 48 8b 4b 08                 	movq	0x8(%rbx), %rcx
14008b943: 49 39 ce                    	cmpq	%rcx, %r14
14008b946: 0f 84 33 06 00 00           	je	0x14008bf7f <.text+0x8af7f>
14008b94c: 48 39 f0                    	cmpq	%rsi, %rax
14008b94f: 0f 84 48 06 00 00           	je	0x14008bf9d <.text+0x8af9d>
14008b955: 48 89 43 08                 	movq	%rax, 0x8(%rbx)
14008b959: 48 8b 84 24 20 01 00 00     	movq	0x120(%rsp), %rax
14008b961: 48 8b 53 18                 	movq	0x18(%rbx), %rdx
14008b965: 4c 89 43 10                 	movq	%r8, 0x10(%rbx)
14008b969: 48 89 43 18                 	movq	%rax, 0x18(%rbx)
14008b96d: 48 85 c9                    	testq	%rcx, %rcx
14008b970: 74 26                       	je	0x14008b998 <.text+0x8a998>
14008b972: 48 89 94 24 20 01 00 00     	movq	%rdx, 0x120(%rsp)
14008b97a: 45 31 e4                    	xorl	%r12d, %r12d
14008b97d: 66 44 89 21                 	movw	%r12w, (%rcx)
14008b981: 48 39 f1                    	cmpq	%rsi, %rcx
14008b984: 74 12                       	je	0x14008b998 <.text+0x8a998>
14008b986: 48 8b 84 24 20 01 00 00     	movq	0x120(%rsp), %rax
14008b98e: 48 8d 54 00 02              	leaq	0x2(%rax,%rax), %rdx
14008b993: e8 38 e7 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008b998: 48 8b 8c 24 f0 00 00 00     	movq	0xf0(%rsp), %rcx
14008b9a0: 4c 39 f9                    	cmpq	%r15, %rcx
14008b9a3: 0f 84 a5 fd ff ff           	je	0x14008b74e <.text+0x8a74e>
14008b9a9: 48 8b 84 24 00 01 00 00     	movq	0x100(%rsp), %rax
14008b9b1: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008b9b5: e8 16 e7 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008b9ba: e9 8f fd ff ff              	jmp	0x14008b74e <.text+0x8a74e>
14008b9bf: 90                          	nop
14008b9c0: 48 8d 84 24 b0 00 00 00     	leaq	0xb0(%rsp), %rax
14008b9c8: 48 8b 4c 24 38              	movq	0x38(%rsp), %rcx
14008b9cd: 45 31 c0                    	xorl	%r8d, %r8d
14008b9d0: 48 8d 94 24 f0 00 00 00     	leaq	0xf0(%rsp), %rdx
14008b9d8: 48 89 44 24 48              	movq	%rax, 0x48(%rsp)
14008b9dd: 48 8d 84 24 90 00 00 00     	leaq	0x90(%rsp), %rax
14008b9e5: 4c 89 4c 24 50              	movq	%r9, 0x50(%rsp)
14008b9ea: 48 89 44 24 58              	movq	%rax, 0x58(%rsp)
14008b9ef: e8 ec ad 03 00              	callq	0x1400c67e0 <.text+0xc57e0>
14008b9f4: 48 89 84 24 10 01 00 00     	movq	%rax, 0x110(%rsp)
14008b9fc: 48 89 c1                    	movq	%rax, %rcx
14008b9ff: 4c 8b 4c 24 50              	movq	0x50(%rsp), %r9
14008ba04: 48 8b 84 24 f0 00 00 00     	movq	0xf0(%rsp), %rax
14008ba0c: 48 89 84 24 20 01 00 00     	movq	%rax, 0x120(%rsp)
14008ba14: 48 8d 15 61 f3 06 00        	leaq	0x6f361(%rip), %rdx     # 0x1400fad7c
14008ba1b: 4d 89 c8                    	movq	%r9, %r8
14008ba1e: e8 c5 9d 01 00              	callq	0x1400a57e8 <.text+0xa47e8>
14008ba23: 48 8b 94 24 f0 00 00 00     	movq	0xf0(%rsp), %rdx
14008ba2b: 48 8b 8c 24 10 01 00 00     	movq	0x110(%rsp), %rcx
14008ba33: 48 8d 04 12                 	leaq	(%rdx,%rdx), %rax
14008ba37: e9 61 fc ff ff              	jmp	0x14008b69d <.text+0x8a69d>
14008ba3c: 0f 1f 40 00                 	nopl	(%rax)
14008ba40: b8 2f 00 00 00              	movl	$0x2f, %eax
14008ba45: 66 89 01                    	movw	%ax, (%rcx)
14008ba48: b8 02 00 00 00              	movl	$0x2, %eax
14008ba4d: e9 4b fc ff ff              	jmp	0x14008b69d <.text+0x8a69d>
14008ba52: 48 8b 48 08                 	movq	0x8(%rax), %rcx
14008ba56: 4c 8d 41 01                 	leaq	0x1(%rcx), %r8
14008ba5a: 41 83 e0 3f                 	andl	$0x3f, %r8d
14008ba5e: 0f 84 31 fb ff ff           	je	0x14008b595 <.text+0x8a595>
14008ba64: 31 c9                       	xorl	%ecx, %ecx
14008ba66: 41 89 c9                    	movl	%ecx, %r9d
14008ba69: 83 c1 01                    	addl	$0x1, %ecx
14008ba6c: 46 0f b6 14 0a              	movzbl	(%rdx,%r9), %r10d
14008ba71: 46 88 94 0c a0 00 00 00     	movb	%r10b, 0xa0(%rsp,%r9)
14008ba79: 44 39 c1                    	cmpl	%r8d, %ecx
14008ba7c: 72 e8                       	jb	0x14008ba66 <.text+0x8aa66>
14008ba7e: e9 0e fb ff ff              	jmp	0x14008b591 <.text+0x8a591>
14008ba83: 48 8b 48 08                 	movq	0x8(%rax), %rcx
14008ba87: 4c 8d 41 01                 	leaq	0x1(%rcx), %r8
14008ba8b: 41 83 e0 3f                 	andl	$0x3f, %r8d
14008ba8f: 0f 84 73 fa ff ff           	je	0x14008b508 <.text+0x8a508>
14008ba95: 31 c9                       	xorl	%ecx, %ecx
14008ba97: 41 89 c9                    	movl	%ecx, %r9d
14008ba9a: 83 c1 01                    	addl	$0x1, %ecx
14008ba9d: 46 0f b6 14 0a              	movzbl	(%rdx,%r9), %r10d
14008baa2: 46 88 14 0e                 	movb	%r10b, (%rsi,%r9)
14008baa6: 44 39 c1                    	cmpl	%r8d, %ecx
14008baa9: 72 ec                       	jb	0x14008ba97 <.text+0x8aa97>
14008baab: e9 54 fa ff ff              	jmp	0x14008b504 <.text+0x8a504>
14008bab0: 48 8d 0d ef f2 06 00        	leaq	0x6f2ef(%rip), %rcx     # 0x1400fada6
14008bab7: 48 89 b4 24 10 01 00 00     	movq	%rsi, 0x110(%rsp)
14008babf: e8 bc 9b 01 00              	callq	0x1400a5680 <.text+0xa4680>
14008bac4: 4c 8d 24 00                 	leaq	(%rax,%rax), %r12
14008bac8: 4d 89 e0                    	movq	%r12, %r8
14008bacb: 49 d1 f8                    	sarq	%r8
14008bace: 4c 89 44 24 68              	movq	%r8, 0x68(%rsp)
14008bad3: 49 83 fc 0e                 	cmpq	$0xe, %r12
14008bad7: 0f 87 17 05 00 00           	ja	0x14008bff4 <.text+0x8aff4>
14008badd: 48 8b 84 24 10 01 00 00     	movq	0x110(%rsp), %rax
14008bae5: 49 83 f8 01                 	cmpq	$0x1, %r8
14008bae9: 0f 84 f1 04 00 00           	je	0x14008bfe0 <.text+0x8afe0>
14008baef: 31 d2                       	xorl	%edx, %edx
14008baf1: 4d 85 c0                    	testq	%r8, %r8
14008baf4: 0f 85 21 05 00 00           	jne	0x14008c01b <.text+0x8b01b>
14008bafa: 31 c9                       	xorl	%ecx, %ecx
14008bafc: 66 89 0c 10                 	movw	%cx, (%rax,%rdx)
14008bb00: 4c 89 84 24 18 01 00 00     	movq	%r8, 0x118(%rsp)
14008bb08: e9 32 fe ff ff              	jmp	0x14008b93f <.text+0x8a93f>
14008bb0d: 81 39 74 72 75 65           	cmpl	$0x65757274, (%rcx)     # imm = 0x65757274
14008bb13: 0f 85 c6 fd ff ff           	jne	0x14008b8df <.text+0x8a8df>
14008bb19: 48 39 f1                    	cmpq	%rsi, %rcx
14008bb1c: 74 11                       	je	0x14008bb2f <.text+0x8ab2f>
14008bb1e: 48 8b 84 24 20 01 00 00     	movq	0x120(%rsp), %rax
14008bb26: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008bb2a: e8 a1 e5 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008bb2f: c6 03 01                    	movb	$0x1, (%rbx)
14008bb32: 48 8b 54 24 48              	movq	0x48(%rsp), %rdx
14008bb37: 4c 8d 05 72 f2 06 00        	leaq	0x6f272(%rip), %r8      # 0x1400fadb0
14008bb3e: 48 8b 4c 24 38              	movq	0x38(%rsp), %rcx
14008bb43: e8 f8 d7 ff ff              	callq	0x140089340 <.text+0x88340>
14008bb48: 48 8b 8c 24 10 01 00 00     	movq	0x110(%rsp), %rcx
14008bb50: 41 b8 0a 00 00 00           	movl	$0xa, %r8d
14008bb56: 31 d2                       	xorl	%edx, %edx
14008bb58: ff 15 12 a5 0e 00           	callq	*0xea512(%rip)          # 0x140176070
14008bb5e: 48 8b 8c 24 10 01 00 00     	movq	0x110(%rsp), %rcx
14008bb66: 48 89 43 28                 	movq	%rax, 0x28(%rbx)
14008bb6a: 48 39 f1                    	cmpq	%rsi, %rcx
14008bb6d: 74 11                       	je	0x14008bb80 <.text+0x8ab80>
14008bb6f: 48 8b 84 24 20 01 00 00     	movq	0x120(%rsp), %rax
14008bb77: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008bb7b: e8 50 e5 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008bb80: 48 8b 54 24 48              	movq	0x48(%rsp), %rdx
14008bb85: 48 8b 4c 24 38              	movq	0x38(%rsp), %rcx
14008bb8a: 4c 8d 05 22 f2 06 00        	leaq	0x6f222(%rip), %r8      # 0x1400fadb3
14008bb91: e8 aa d7 ff ff              	callq	0x140089340 <.text+0x88340>
14008bb96: 48 8b 8c 24 10 01 00 00     	movq	0x110(%rsp), %rcx
14008bb9e: 41 b8 0a 00 00 00           	movl	$0xa, %r8d
14008bba4: 31 d2                       	xorl	%edx, %edx
14008bba6: ff 15 bc a4 0e 00           	callq	*0xea4bc(%rip)          # 0x140176068
14008bbac: 48 8b 8c 24 10 01 00 00     	movq	0x110(%rsp), %rcx
14008bbb4: 48 89 43 30                 	movq	%rax, 0x30(%rbx)
14008bbb8: 48 39 f1                    	cmpq	%rsi, %rcx
14008bbbb: 74 11                       	je	0x14008bbce <.text+0x8abce>
14008bbbd: 48 8b 84 24 20 01 00 00     	movq	0x120(%rsp), %rax
14008bbc5: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008bbc9: e8 02 e5 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008bbce: 48 8b 54 24 48              	movq	0x48(%rsp), %rdx
14008bbd3: 48 8b 4c 24 38              	movq	0x38(%rsp), %rcx
14008bbd8: 4c 8d 05 df f1 06 00        	leaq	0x6f1df(%rip), %r8      # 0x1400fadbe
14008bbdf: e8 5c d7 ff ff              	callq	0x140089340 <.text+0x88340>
14008bbe4: 48 8b 8c 24 10 01 00 00     	movq	0x110(%rsp), %rcx
14008bbec: e8 af 9c 01 00              	callq	0x1400a58a0 <.text+0xa48a0>
14008bbf1: 48 8b 8c 24 10 01 00 00     	movq	0x110(%rsp), %rcx
14008bbf9: 89 43 38                    	movl	%eax, 0x38(%rbx)
14008bbfc: 48 39 f1                    	cmpq	%rsi, %rcx
14008bbff: 74 11                       	je	0x14008bc12 <.text+0x8ac12>
14008bc01: 48 8b 84 24 20 01 00 00     	movq	0x120(%rsp), %rax
14008bc09: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008bc0d: e8 be e4 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008bc12: 48 8b 54 24 48              	movq	0x48(%rsp), %rdx
14008bc17: 48 8b 4c 24 38              	movq	0x38(%rsp), %rcx
14008bc1c: 4c 8d 05 a7 f1 06 00        	leaq	0x6f1a7(%rip), %r8      # 0x1400fadca
14008bc23: e8 18 d7 ff ff              	callq	0x140089340 <.text+0x88340>
14008bc28: 48 8b 8c 24 10 01 00 00     	movq	0x110(%rsp), %rcx
14008bc30: 48 83 bc 24 18 01 00 00 04  	cmpq	$0x4, 0x118(%rsp)
14008bc39: 0f 84 0c 02 00 00           	je	0x14008be4b <.text+0x8ae4b>
14008bc3f: c6 43 3c 00                 	movb	$0x0, 0x3c(%rbx)
14008bc43: 48 39 f1                    	cmpq	%rsi, %rcx
14008bc46: 74 11                       	je	0x14008bc59 <.text+0x8ac59>
14008bc48: 48 8b 84 24 20 01 00 00     	movq	0x120(%rsp), %rax
14008bc50: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008bc54: e8 77 e4 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008bc59: 48 8b 54 24 48              	movq	0x48(%rsp), %rdx
14008bc5e: 4c 8d 05 6f f1 06 00        	leaq	0x6f16f(%rip), %r8      # 0x1400fadd4
14008bc65: 48 8d 8c 24 f0 00 00 00     	leaq	0xf0(%rsp), %rcx
14008bc6d: e8 ce d6 ff ff              	callq	0x140089340 <.text+0x88340>
14008bc72: 48 8b 4c 24 38              	movq	0x38(%rsp), %rcx
14008bc77: 48 8d 94 24 f0 00 00 00     	leaq	0xf0(%rsp), %rdx
14008bc7f: e8 4c d4 ff ff              	callq	0x1400890d0 <.text+0x880d0>
14008bc84: 48 8b 4b 40                 	movq	0x40(%rbx), %rcx
14008bc88: 4c 8b 84 24 18 01 00 00     	movq	0x118(%rsp), %r8
14008bc90: 48 8b 84 24 10 01 00 00     	movq	0x110(%rsp), %rax
14008bc98: 49 39 cd                    	cmpq	%rcx, %r13
14008bc9b: 0f 84 bf 02 00 00           	je	0x14008bf60 <.text+0x8af60>
14008bca1: 48 39 f0                    	cmpq	%rsi, %rax
14008bca4: 0f 84 74 02 00 00           	je	0x14008bf1e <.text+0x8af1e>
14008bcaa: 48 8b 53 50                 	movq	0x50(%rbx), %rdx
14008bcae: 66 49 0f 6e c0              	movq	%r8, %xmm0
14008bcb3: 48 89 43 40                 	movq	%rax, 0x40(%rbx)
14008bcb7: 0f 16 84 24 20 01 00 00     	movhps	0x120(%rsp), %xmm0      # xmm0 = xmm0[0,1],mem[0,1]
14008bcbf: 0f 11 43 48                 	movups	%xmm0, 0x48(%rbx)
14008bcc3: 48 85 c9                    	testq	%rcx, %rcx
14008bcc6: 74 26                       	je	0x14008bcee <.text+0x8acee>
14008bcc8: 48 89 94 24 20 01 00 00     	movq	%rdx, 0x120(%rsp)
14008bcd0: 45 31 d2                    	xorl	%r10d, %r10d
14008bcd3: 66 44 89 11                 	movw	%r10w, (%rcx)
14008bcd7: 48 39 f1                    	cmpq	%rsi, %rcx
14008bcda: 74 12                       	je	0x14008bcee <.text+0x8acee>
14008bcdc: 48 8b 84 24 20 01 00 00     	movq	0x120(%rsp), %rax
14008bce4: 48 8d 54 00 02              	leaq	0x2(%rax,%rax), %rdx
14008bce9: e8 e2 e3 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008bcee: 48 8b 8c 24 f0 00 00 00     	movq	0xf0(%rsp), %rcx
14008bcf6: 4c 39 f9                    	cmpq	%r15, %rcx
14008bcf9: 74 11                       	je	0x14008bd0c <.text+0x8ad0c>
14008bcfb: 48 8b 84 24 00 01 00 00     	movq	0x100(%rsp), %rax
14008bd03: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008bd07: e8 c4 e3 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008bd0c: 48 8b 54 24 48              	movq	0x48(%rsp), %rdx
14008bd11: 4c 8d 05 c1 f0 06 00        	leaq	0x6f0c1(%rip), %r8      # 0x1400fadd9
14008bd18: 48 8d 8c 24 f0 00 00 00     	leaq	0xf0(%rsp), %rcx
14008bd20: e8 1b d6 ff ff              	callq	0x140089340 <.text+0x88340>
14008bd25: 48 8b 4c 24 38              	movq	0x38(%rsp), %rcx
14008bd2a: 48 8d 94 24 f0 00 00 00     	leaq	0xf0(%rsp), %rdx
14008bd32: e8 99 d3 ff ff              	callq	0x1400890d0 <.text+0x880d0>
14008bd37: 48 8b 4b 60                 	movq	0x60(%rbx), %rcx
14008bd3b: 48 8b 84 24 10 01 00 00     	movq	0x110(%rsp), %rax
14008bd43: 4c 8b 84 24 18 01 00 00     	movq	0x118(%rsp), %r8
14008bd4b: 49 39 cc                    	cmpq	%rcx, %r12
14008bd4e: 0f 84 ab 01 00 00           	je	0x14008beff <.text+0x8aeff>
14008bd54: 48 39 f0                    	cmpq	%rsi, %rax
14008bd57: 0f 84 60 01 00 00           	je	0x14008bebd <.text+0x8aebd>
14008bd5d: 48 8b 53 70                 	movq	0x70(%rbx), %rdx
14008bd61: 66 49 0f 6e c0              	movq	%r8, %xmm0
14008bd66: 48 89 43 60                 	movq	%rax, 0x60(%rbx)
14008bd6a: 0f 16 84 24 20 01 00 00     	movhps	0x120(%rsp), %xmm0      # xmm0 = xmm0[0,1],mem[0,1]
14008bd72: 0f 11 43 68                 	movups	%xmm0, 0x68(%rbx)
14008bd76: 48 85 c9                    	testq	%rcx, %rcx
14008bd79: 74 24                       	je	0x14008bd9f <.text+0x8ad9f>
14008bd7b: 48 89 94 24 20 01 00 00     	movq	%rdx, 0x120(%rsp)
14008bd83: 31 d2                       	xorl	%edx, %edx
14008bd85: 66 89 11                    	movw	%dx, (%rcx)
14008bd88: 48 39 f1                    	cmpq	%rsi, %rcx
14008bd8b: 74 12                       	je	0x14008bd9f <.text+0x8ad9f>
14008bd8d: 48 8b 84 24 20 01 00 00     	movq	0x120(%rsp), %rax
14008bd95: 48 8d 54 00 02              	leaq	0x2(%rax,%rax), %rdx
14008bd9a: e8 31 e3 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008bd9f: 48 8b 8c 24 f0 00 00 00     	movq	0xf0(%rsp), %rcx
14008bda7: 4c 39 f9                    	cmpq	%r15, %rcx
14008bdaa: 74 11                       	je	0x14008bdbd <.text+0x8adbd>
14008bdac: 48 8b 84 24 00 01 00 00     	movq	0x100(%rsp), %rax
14008bdb4: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008bdb8: e8 13 e3 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008bdbd: 48 8b 73 68                 	movq	0x68(%rbx), %rsi
14008bdc1: 48 8d 0d 6c ef 06 00        	leaq	0x6ef6c(%rip), %rcx     # 0x1400fad34
14008bdc8: e8 b3 98 01 00              	callq	0x1400a5680 <.text+0xa4680>
14008bdcd: 48 39 c6                    	cmpq	%rax, %rsi
14008bdd0: 0f 84 93 00 00 00           	je	0x14008be69 <.text+0x8ae69>
14008bdd6: 4c 8b 63 68                 	movq	0x68(%rbx), %r12
14008bdda: 48 8d 35 6b ef 06 00        	leaq	0x6ef6b(%rip), %rsi     # 0x1400fad4c
14008bde1: 48 89 f1                    	movq	%rsi, %rcx
14008bde4: e8 97 98 01 00              	callq	0x1400a5680 <.text+0xa4680>
14008bde9: 49 39 c4                    	cmpq	%rax, %r12
14008bdec: 0f 84 67 02 00 00           	je	0x14008c059 <.text+0x8b059>
14008bdf2: 48 8d 0d e5 ef 06 00        	leaq	0x6efe5(%rip), %rcx     # 0x1400fadde
14008bdf9: e8 82 98 01 00              	callq	0x1400a5680 <.text+0xa4680>
14008bdfe: 48 8d 4b 08                 	leaq	0x8(%rbx), %rcx
14008be02: 4c 8d 0d d5 ef 06 00        	leaq	0x6efd5(%rip), %r9      # 0x1400fadde
14008be09: 31 d2                       	xorl	%edx, %edx
14008be0b: 48 89 44 24 20              	movq	%rax, 0x20(%rsp)
14008be10: 4c 8b 43 10                 	movq	0x10(%rbx), %r8
14008be14: e8 07 7f 03 00              	callq	0x1400c3d20 <.text+0xc2d20>
14008be19: e9 30 f9 ff ff              	jmp	0x14008b74e <.text+0x8a74e>
14008be1e: 48 8b 48 08                 	movq	0x8(%rax), %rcx
14008be22: 4c 8d 41 01                 	leaq	0x1(%rcx), %r8
14008be26: 41 83 e0 3f                 	andl	$0x3f, %r8d
14008be2a: 0f 84 59 fa ff ff           	je	0x14008b889 <.text+0x8a889>
14008be30: 31 c9                       	xorl	%ecx, %ecx
14008be32: 41 89 c9                    	movl	%ecx, %r9d
14008be35: 83 c1 01                    	addl	$0x1, %ecx
14008be38: 46 0f b6 14 0a              	movzbl	(%rdx,%r9), %r10d
14008be3d: 46 88 14 0e                 	movb	%r10b, (%rsi,%r9)
14008be41: 44 39 c1                    	cmpl	%r8d, %ecx
14008be44: 72 ec                       	jb	0x14008be32 <.text+0x8ae32>
14008be46: e9 3a fa ff ff              	jmp	0x14008b885 <.text+0x8a885>
14008be4b: 81 39 74 72 75 65           	cmpl	$0x65757274, (%rcx)     # imm = 0x65757274
14008be51: 0f 85 e8 fd ff ff           	jne	0x14008bc3f <.text+0x8ac3f>
14008be57: c6 43 3c 01                 	movb	$0x1, 0x3c(%rbx)
14008be5b: 48 39 f1                    	cmpq	%rsi, %rcx
14008be5e: 0f 85 e4 fd ff ff           	jne	0x14008bc48 <.text+0x8ac48>
14008be64: e9 f0 fd ff ff              	jmp	0x14008bc59 <.text+0x8ac59>
14008be69: 48 8b 53 60                 	movq	0x60(%rbx), %rdx
14008be6d: 48 8d 05 c0 ee 06 00        	leaq	0x6eec0(%rip), %rax     # 0x1400fad34
14008be74: 4c 8b 63 68                 	movq	0x68(%rbx), %r12
14008be78: 48 39 c2                    	cmpq	%rax, %rdx
14008be7b: 74 2e                       	je	0x14008beab <.text+0x8aeab>
14008be7d: 4d 85 e4                    	testq	%r12, %r12
14008be80: 74 29                       	je	0x14008beab <.text+0x8aeab>
14008be82: 48 85 d2                    	testq	%rdx, %rdx
14008be85: 0f 84 4f ff ff ff           	je	0x14008bdda <.text+0x8adda>
14008be8b: 4d 89 e0                    	movq	%r12, %r8
14008be8e: 48 89 d1                    	movq	%rdx, %rcx
14008be91: 0f b7 30                    	movzwl	(%rax), %esi
14008be94: 66 39 31                    	cmpw	%si, (%rcx)
14008be97: 0f 85 3d ff ff ff           	jne	0x14008bdda <.text+0x8adda>
14008be9d: 48 83 c1 02                 	addq	$0x2, %rcx
14008bea1: 48 83 c0 02                 	addq	$0x2, %rax
14008bea5: 49 83 e8 01                 	subq	$0x1, %r8
14008bea9: 75 e6                       	jne	0x14008be91 <.text+0x8ae91>
14008beab: 31 c0                       	xorl	%eax, %eax
14008bead: 66 89 02                    	movw	%ax, (%rdx)
14008beb0: 48 c7 43 68 00 00 00 00     	movq	$0x0, 0x68(%rbx)
14008beb8: e9 35 ff ff ff              	jmp	0x14008bdf2 <.text+0x8adf2>
14008bebd: 4d 85 c0                    	testq	%r8, %r8
14008bec0: 74 29                       	je	0x14008beeb <.text+0x8aeeb>
14008bec2: 49 83 f8 01                 	cmpq	$0x1, %r8
14008bec6: 0f 84 e4 01 00 00           	je	0x14008c0b0 <.text+0x8b0b0>
14008becc: 4d 01 c0                    	addq	%r8, %r8
14008becf: 48 89 f2                    	movq	%rsi, %rdx
14008bed2: e8 11 99 01 00              	callq	0x1400a57e8 <.text+0xa47e8>
14008bed7: 4c 8b 84 24 18 01 00 00     	movq	0x118(%rsp), %r8
14008bedf: 48 8b 4b 60                 	movq	0x60(%rbx), %rcx
14008bee3: 48 8b 84 24 10 01 00 00     	movq	0x110(%rsp), %rax
14008beeb: 45 31 c9                    	xorl	%r9d, %r9d
14008beee: 66 46 89 0c 41              	movw	%r9w, (%rcx,%r8,2)
14008bef3: 48 89 c1                    	movq	%rax, %rcx
14008bef6: 4c 89 43 68                 	movq	%r8, 0x68(%rbx)
14008befa: e9 84 fe ff ff              	jmp	0x14008bd83 <.text+0x8ad83>
14008beff: 48 39 f0                    	cmpq	%rsi, %rax
14008bf02: 74 b9                       	je	0x14008bebd <.text+0x8aebd>
14008bf04: 66 49 0f 6e c0              	movq	%r8, %xmm0
14008bf09: 48 89 43 60                 	movq	%rax, 0x60(%rbx)
14008bf0d: 0f 16 84 24 20 01 00 00     	movhps	0x120(%rsp), %xmm0      # xmm0 = xmm0[0,1],mem[0,1]
14008bf15: 0f 11 43 68                 	movups	%xmm0, 0x68(%rbx)
14008bf19: e9 81 fe ff ff              	jmp	0x14008bd9f <.text+0x8ad9f>
14008bf1e: 4d 85 c0                    	testq	%r8, %r8
14008bf21: 74 29                       	je	0x14008bf4c <.text+0x8af4c>
14008bf23: 49 83 f8 01                 	cmpq	$0x1, %r8
14008bf27: 0f 84 73 01 00 00           	je	0x14008c0a0 <.text+0x8b0a0>
14008bf2d: 4d 01 c0                    	addq	%r8, %r8
14008bf30: 48 89 f2                    	movq	%rsi, %rdx
14008bf33: e8 b0 98 01 00              	callq	0x1400a57e8 <.text+0xa47e8>
14008bf38: 4c 8b 84 24 18 01 00 00     	movq	0x118(%rsp), %r8
14008bf40: 48 8b 4b 40                 	movq	0x40(%rbx), %rcx
14008bf44: 48 8b 84 24 10 01 00 00     	movq	0x110(%rsp), %rax
14008bf4c: 45 31 db                    	xorl	%r11d, %r11d
14008bf4f: 66 46 89 1c 41              	movw	%r11w, (%rcx,%r8,2)
14008bf54: 48 89 c1                    	movq	%rax, %rcx
14008bf57: 4c 89 43 48                 	movq	%r8, 0x48(%rbx)
14008bf5b: e9 70 fd ff ff              	jmp	0x14008bcd0 <.text+0x8acd0>
14008bf60: 48 39 f0                    	cmpq	%rsi, %rax
14008bf63: 74 b9                       	je	0x14008bf1e <.text+0x8af1e>
14008bf65: 66 49 0f 6e c0              	movq	%r8, %xmm0
14008bf6a: 48 89 43 40                 	movq	%rax, 0x40(%rbx)
14008bf6e: 0f 16 84 24 20 01 00 00     	movhps	0x120(%rsp), %xmm0      # xmm0 = xmm0[0,1],mem[0,1]
14008bf76: 0f 11 43 48                 	movups	%xmm0, 0x48(%rbx)
14008bf7a: e9 6f fd ff ff              	jmp	0x14008bcee <.text+0x8acee>
14008bf7f: 48 39 f0                    	cmpq	%rsi, %rax
14008bf82: 74 19                       	je	0x14008bf9d <.text+0x8af9d>
14008bf84: 48 89 43 08                 	movq	%rax, 0x8(%rbx)
14008bf88: 48 8b 84 24 20 01 00 00     	movq	0x120(%rsp), %rax
14008bf90: 4c 89 43 10                 	movq	%r8, 0x10(%rbx)
14008bf94: 48 89 43 18                 	movq	%rax, 0x18(%rbx)
14008bf98: e9 fb f9 ff ff              	jmp	0x14008b998 <.text+0x8a998>
14008bf9d: 4d 85 c0                    	testq	%r8, %r8
14008bfa0: 0f 84 a8 00 00 00           	je	0x14008c04e <.text+0x8b04e>
14008bfa6: 49 83 f8 01                 	cmpq	$0x1, %r8
14008bfaa: 0f 84 93 00 00 00           	je	0x14008c043 <.text+0x8b043>
14008bfb0: 4d 01 c0                    	addq	%r8, %r8
14008bfb3: 48 89 f2                    	movq	%rsi, %rdx
14008bfb6: e8 2d 98 01 00              	callq	0x1400a57e8 <.text+0xa47e8>
14008bfbb: 4c 8b 84 24 18 01 00 00     	movq	0x118(%rsp), %r8
14008bfc3: 48 8b 53 08                 	movq	0x8(%rbx), %rdx
14008bfc7: 48 8b 8c 24 10 01 00 00     	movq	0x110(%rsp), %rcx
14008bfcf: 45 31 ed                    	xorl	%r13d, %r13d
14008bfd2: 66 46 89 2c 42              	movw	%r13w, (%rdx,%r8,2)
14008bfd7: 4c 89 43 10                 	movq	%r8, 0x10(%rbx)
14008bfdb: e9 9a f9 ff ff              	jmp	0x14008b97a <.text+0x8a97a>
14008bfe0: 41 b9 61 53 00 00           	movl	$0x5361, %r9d           # imm = 0x5361
14008bfe6: ba 02 00 00 00              	movl	$0x2, %edx
14008bfeb: 66 44 89 08                 	movw	%r9w, (%rax)
14008bfef: e9 06 fb ff ff              	jmp	0x14008bafa <.text+0x8aafa>
14008bff4: 48 8b 4c 24 38              	movq	0x38(%rsp), %rcx
14008bff9: 48 8d 54 24 68              	leaq	0x68(%rsp), %rdx
14008bffe: 45 31 c0                    	xorl	%r8d, %r8d
14008c001: e8 da a7 03 00              	callq	0x1400c67e0 <.text+0xc57e0>
14008c006: 48 8b 54 24 68              	movq	0x68(%rsp), %rdx
14008c00b: 48 89 84 24 10 01 00 00     	movq	%rax, 0x110(%rsp)
14008c013: 48 89 94 24 20 01 00 00     	movq	%rdx, 0x120(%rsp)
14008c01b: 4d 89 e0                    	movq	%r12, %r8
14008c01e: 48 8d 15 81 ed 06 00        	leaq	0x6ed81(%rip), %rdx     # 0x1400fada6
14008c025: 48 89 c1                    	movq	%rax, %rcx
14008c028: e8 bb 97 01 00              	callq	0x1400a57e8 <.text+0xa47e8>
14008c02d: 4c 8b 44 24 68              	movq	0x68(%rsp), %r8
14008c032: 48 8b 84 24 10 01 00 00     	movq	0x110(%rsp), %rax
14008c03a: 4b 8d 14 00                 	leaq	(%r8,%r8), %rdx
14008c03e: e9 b7 fa ff ff              	jmp	0x14008bafa <.text+0x8aafa>
14008c043: 0f b7 94 24 20 01 00 00     	movzwl	0x120(%rsp), %edx
14008c04b: 66 89 11                    	movw	%dx, (%rcx)
14008c04e: 48 89 ca                    	movq	%rcx, %rdx
14008c051: 48 89 c1                    	movq	%rax, %rcx
14008c054: e9 76 ff ff ff              	jmp	0x14008bfcf <.text+0x8afcf>
14008c059: 48 8b 43 68                 	movq	0x68(%rbx), %rax
14008c05d: 48 8b 53 60                 	movq	0x60(%rbx), %rdx
14008c061: 48 85 c0                    	testq	%rax, %rax
14008c064: 0f 84 41 fe ff ff           	je	0x14008beab <.text+0x8aeab>
14008c06a: 48 39 f2                    	cmpq	%rsi, %rdx
14008c06d: 0f 84 38 fe ff ff           	je	0x14008beab <.text+0x8aeab>
14008c073: 48 85 d2                    	testq	%rdx, %rdx
14008c076: 0f 84 76 fd ff ff           	je	0x14008bdf2 <.text+0x8adf2>
14008c07c: 48 89 d1                    	movq	%rdx, %rcx
14008c07f: 44 0f b7 1e                 	movzwl	(%rsi), %r11d
14008c083: 66 44 39 19                 	cmpw	%r11w, (%rcx)
14008c087: 0f 85 65 fd ff ff           	jne	0x14008bdf2 <.text+0x8adf2>
14008c08d: 48 83 c1 02                 	addq	$0x2, %rcx
14008c091: 48 83 c6 02                 	addq	$0x2, %rsi
14008c095: 48 83 e8 01                 	subq	$0x1, %rax
14008c099: 75 e4                       	jne	0x14008c07f <.text+0x8b07f>
14008c09b: e9 0b fe ff ff              	jmp	0x14008beab <.text+0x8aeab>
14008c0a0: 0f b7 94 24 20 01 00 00     	movzwl	0x120(%rsp), %edx
14008c0a8: 66 89 11                    	movw	%dx, (%rcx)
14008c0ab: e9 9c fe ff ff              	jmp	0x14008bf4c <.text+0x8af4c>
14008c0b0: 0f b7 94 24 20 01 00 00     	movzwl	0x120(%rsp), %edx
14008c0b8: 66 89 11                    	movw	%dx, (%rcx)
14008c0bb: e9 2b fe ff ff              	jmp	0x14008beeb <.text+0x8aeeb>
14008c0c0: 48 89 c6                    	movq	%rax, %rsi
14008c0c3: e9 a8 c2 04 00              	jmp	0x1400d8370 <.text+0xd7370>
14008c0c8: 48 89 c6                    	movq	%rax, %rsi
14008c0cb: e9 e5 c2 04 00              	jmp	0x1400d83b5 <.text+0xd73b5>
14008c0d0: 48 89 c6                    	movq	%rax, %rsi
14008c0d3: e9 05 c3 04 00              	jmp	0x1400d83dd <.text+0xd73dd>
14008c0d8: 48 89 c6                    	movq	%rax, %rsi
14008c0db: e9 28 c3 04 00              	jmp	0x1400d8408 <.text+0xd7408>
14008c0e0: 48 89 c6                    	movq	%rax, %rsi
14008c0e3: e9 2f c3 04 00              	jmp	0x1400d8417 <.text+0xd7417>
14008c0e8: 48 89 c6                    	movq	%rax, %rsi
14008c0eb: e9 ab c2 04 00              	jmp	0x1400d839b <.text+0xd739b>
14008c0f0: 48 89 c6                    	movq	%rax, %rsi
14008c0f3: e9 85 c2 04 00              	jmp	0x1400d837d <.text+0xd737d>
14008c0f8: 48 89 c6                    	movq	%rax, %rsi
14008c0fb: 48 8d 84 24 b0 00 00 00     	leaq	0xb0(%rsp), %rax
14008c103: 48 89 44 24 48              	movq	%rax, 0x48(%rsp)
14008c108: 48 8d 84 24 90 00 00 00     	leaq	0x90(%rsp), %rax
14008c110: 48 89 44 24 58              	movq	%rax, 0x58(%rsp)
14008c115: e9 19 c3 04 00              	jmp	0x1400d8433 <.text+0xd7433>
14008c11a: 48 89 c6                    	movq	%rax, %rsi
14008c11d: e9 04 c3 04 00              	jmp	0x1400d8426 <.text+0xd7426>
14008c122: e9 34 c3 04 00              	jmp	0x1400d845b <.text+0xd745b>
14008c127: e9 c0 c2 04 00              	jmp	0x1400d83ec <.text+0xd73ec>
14008c12c: 48 89 c6                    	movq	%rax, %rsi
14008c12f: e9 71 c2 04 00              	jmp	0x1400d83a5 <.text+0xd73a5>
14008c134: 90                          	nop
14008c135: 66 66 2e 0f 1f 84 00 00 00 00 00    	nopw	%cs:(%rax,%rax)
14008c140: 41 57                       	pushq	%r15
14008c142: 41 56                       	pushq	%r14
14008c144: 41 55                       	pushq	%r13
14008c146: 41 54                       	pushq	%r12
14008c148: 55                          	pushq	%rbp
14008c149: 57                          	pushq	%rdi
14008c14a: 56                          	pushq	%rsi
14008c14b: 53                          	pushq	%rbx
14008c14c: 48 81 ec 88 02 00 00        	subq	$0x288, %rsp            # imm = 0x288
14008c153: 31 f6                       	xorl	%esi, %esi
14008c155: 45 31 e4                    	xorl	%r12d, %r12d
14008c158: 31 ed                       	xorl	%ebp, %ebp
14008c15a: 45 31 db                    	xorl	%r11d, %r11d
14008c15d: 45 31 ff                    	xorl	%r15d, %r15d
14008c160: 66 89 71 18                 	movw	%si, 0x18(%rcx)
14008c164: 48 8d 79 18                 	leaq	0x18(%rcx), %rdi
14008c168: 4c 8d 71 50                 	leaq	0x50(%rcx), %r14
14008c16c: 48 89 cb                    	movq	%rcx, %rbx
14008c16f: 66 44 89 61 70              	movw	%r12w, 0x70(%rcx)
14008c174: 4c 8d 69 70                 	leaq	0x70(%rcx), %r13
14008c178: 4c 8d a1 90 00 00 00        	leaq	0x90(%rcx), %r12
14008c17f: 48 8d b1 b0 00 00 00        	leaq	0xb0(%rcx), %rsi
14008c186: c6 81 b0 00 00 00 00        	movb	$0x0, 0xb0(%rcx)
14008c18d: 48 89 79 08                 	movq	%rdi, 0x8(%rcx)
14008c191: 48 c7 41 10 00 00 00 00     	movq	$0x0, 0x10(%rcx)
14008c199: 48 c7 41 28 00 00 00 00     	movq	$0x0, 0x28(%rcx)
14008c1a1: 48 c7 41 30 00 00 00 00     	movq	$0x0, 0x30(%rcx)
14008c1a9: 4c 89 71 40                 	movq	%r14, 0x40(%rcx)
14008c1ad: 66 89 69 50                 	movw	%bp, 0x50(%rcx)
14008c1b1: 48 8d 6c 24 70              	leaq	0x70(%rsp), %rbp
14008c1b6: 48 c7 41 48 00 00 00 00     	movq	$0x0, 0x48(%rcx)
14008c1be: 4c 89 69 60                 	movq	%r13, 0x60(%rcx)
14008c1c2: 48 c7 41 68 00 00 00 00     	movq	$0x0, 0x68(%rcx)
14008c1ca: 4c 89 a1 80 00 00 00        	movq	%r12, 0x80(%rcx)
14008c1d1: 48 8b 02                    	movq	(%rdx), %rax
14008c1d4: 48 c7 81 88 00 00 00 00 00 00 00    	movq	$0x0, 0x88(%rcx)
14008c1df: 48 c7 81 a8 00 00 00 00 00 00 00    	movq	$0x0, 0xa8(%rcx)
14008c1ea: 4c 8b 42 08                 	movq	0x8(%rdx), %r8
14008c1ee: 48 89 c2                    	movq	%rax, %rdx
14008c1f1: 66 44 89 19                 	movw	%r11w, (%rcx)
14008c1f5: c7 41 38 00 00 00 00        	movl	$0x0, 0x38(%rcx)
14008c1fc: c6 41 3c 00                 	movb	$0x0, 0x3c(%rcx)
14008c200: 66 44 89 b9 90 00 00 00     	movw	%r15w, 0x90(%rcx)
14008c208: 48 89 b1 a0 00 00 00        	movq	%rsi, 0xa0(%rcx)
14008c20f: 48 8d 4c 24 60              	leaq	0x60(%rsp), %rcx
14008c214: 48 89 6c 24 60              	movq	%rbp, 0x60(%rsp)
14008c219: e8 f2 7e 03 00              	callq	0x1400c4110 <.text+0xc3110>
14008c21e: 48 8b 44 24 60              	movq	0x60(%rsp), %rax
14008c223: 48 8b 54 24 68              	movq	0x68(%rsp), %rdx
14008c228: 4c 8d 0c 50                 	leaq	(%rax,%rdx,2), %r9
14008c22c: 4c 39 c8                    	cmpq	%r9, %rax
14008c22f: 74 6a                       	je	0x14008c29b <.text+0x8b29b>
14008c231: 0f 1f 40 00                 	nopl	(%rax)
14008c235: 66 66 2e 0f 1f 84 00 00 00 00 00    	nopw	%cs:(%rax,%rax)
14008c240: 0f b7 08                    	movzwl	(%rax), %ecx
14008c243: 44 8d 41 9f                 	leal	-0x61(%rcx), %r8d
14008c247: 66 41 83 f8 19              	cmpw	$0x19, %r8w
14008c24c: 77 06                       	ja	0x14008c254 <.text+0x8b254>
14008c24e: 83 e9 20                    	subl	$0x20, %ecx
14008c251: 66 89 08                    	movw	%cx, (%rax)
14008c254: 48 83 c0 02                 	addq	$0x2, %rax
14008c258: 49 39 c1                    	cmpq	%rax, %r9
14008c25b: 75 e3                       	jne	0x14008c240 <.text+0x8b240>
14008c25d: 48 85 d2                    	testq	%rdx, %rdx
14008c260: 74 3e                       	je	0x14008c2a0 <.text+0x8b2a0>
14008c262: 0f 1f 00                    	nopl	(%rax)
14008c265: 66 66 2e 0f 1f 84 00 00 00 00 00    	nopw	%cs:(%rax,%rax)
14008c270: 48 8b 44 24 60              	movq	0x60(%rsp), %rax
14008c275: 0f b7 08                    	movzwl	(%rax), %ecx
14008c278: 66 83 f9 20                 	cmpw	$0x20, %cx
14008c27c: 74 06                       	je	0x14008c284 <.text+0x8b284>
14008c27e: 66 83 f9 09                 	cmpw	$0x9, %cx
14008c282: 75 7c                       	jne	0x14008c300 <.text+0x8b300>
14008c284: 41 b8 01 00 00 00           	movl	$0x1, %r8d
14008c28a: 31 d2                       	xorl	%edx, %edx
14008c28c: 48 8d 4c 24 60              	leaq	0x60(%rsp), %rcx
14008c291: e8 ca a2 03 00              	callq	0x1400c6560 <.text+0xc5560>
14008c296: 48 8b 54 24 68              	movq	0x68(%rsp), %rdx
14008c29b: 48 85 d2                    	testq	%rdx, %rdx
14008c29e: 75 d0                       	jne	0x14008c270 <.text+0x8b270>
14008c2a0: 48 8d 0d bf ea 06 00        	leaq	0x6eabf(%rip), %rcx     # 0x1400fad66
14008c2a7: 48 8d 73 08                 	leaq	0x8(%rbx), %rsi
14008c2ab: e8 d0 93 01 00              	callq	0x1400a5680 <.text+0xa4680>
14008c2b0: 4c 8b 43 10                 	movq	0x10(%rbx), %r8
14008c2b4: 31 d2                       	xorl	%edx, %edx
14008c2b6: 48 89 f1                    	movq	%rsi, %rcx
14008c2b9: 48 89 44 24 20              	movq	%rax, 0x20(%rsp)
14008c2be: 4c 8d 0d a1 ea 06 00        	leaq	0x6eaa1(%rip), %r9      # 0x1400fad66
14008c2c5: e8 56 7a 03 00              	callq	0x1400c3d20 <.text+0xc2d20>
14008c2ca: 48 8b 4c 24 60              	movq	0x60(%rsp), %rcx
14008c2cf: 48 39 e9                    	cmpq	%rbp, %rcx
14008c2d2: 74 0f                       	je	0x14008c2e3 <.text+0x8b2e3>
14008c2d4: 48 8b 44 24 70              	movq	0x70(%rsp), %rax
14008c2d9: 48 8d 54 00 02              	leaq	0x2(%rax,%rax), %rdx
14008c2de: e8 ed dd 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008c2e3: 48 89 d8                    	movq	%rbx, %rax
14008c2e6: 48 81 c4 88 02 00 00        	addq	$0x288, %rsp            # imm = 0x288
14008c2ed: 5b                          	popq	%rbx
14008c2ee: 5e                          	popq	%rsi
14008c2ef: 5f                          	popq	%rdi
14008c2f0: 5d                          	popq	%rbp
14008c2f1: 41 5c                       	popq	%r12
14008c2f3: 41 5d                       	popq	%r13
14008c2f5: 41 5e                       	popq	%r14
14008c2f7: 41 5f                       	popq	%r15
14008c2f9: c3                          	retq
14008c2fa: 66 0f 1f 44 00 00           	nopw	(%rax,%rax)
14008c300: 48 89 d1                    	movq	%rdx, %rcx
14008c303: eb 08                       	jmp	0x14008c30d <.text+0x8b30d>
14008c305: 0f 1f 00                    	nopl	(%rax)
14008c308: 48 8b 44 24 60              	movq	0x60(%rsp), %rax
14008c30d: 48 8d 51 ff                 	leaq	-0x1(%rcx), %rdx
14008c311: 44 0f b7 04 50              	movzwl	(%rax,%rdx,2), %r8d
14008c316: 66 41 83 f8 09              	cmpw	$0x9, %r8w
14008c31b: 74 07                       	je	0x14008c324 <.text+0x8b324>
14008c31d: 66 41 83 f8 20              	cmpw	$0x20, %r8w
14008c322: 75 24                       	jne	0x14008c348 <.text+0x8b348>
14008c324: 48 8d 4c 24 60              	leaq	0x60(%rsp), %rcx
14008c329: 41 b8 01 00 00 00           	movl	$0x1, %r8d
14008c32f: e8 2c a2 03 00              	callq	0x1400c6560 <.text+0xc5560>
14008c334: 48 8b 4c 24 68              	movq	0x68(%rsp), %rcx
14008c339: 48 85 c9                    	testq	%rcx, %rcx
14008c33c: 75 ca                       	jne	0x14008c308 <.text+0x8b308>
14008c33e: e9 5d ff ff ff              	jmp	0x14008c2a0 <.text+0x8b2a0>
14008c343: 0f 1f 44 00 00              	nopl	(%rax,%rax)
14008c348: 48 83 f9 10                 	cmpq	$0x10, %rcx
14008c34c: 0f 85 4e ff ff ff           	jne	0x14008c2a0 <.text+0x8b2a0>
14008c352: 48 8d 48 20                 	leaq	0x20(%rax), %rcx
14008c356: 66 2e 0f 1f 84 00 00 00 00 00       	nopw	%cs:(%rax,%rax)
14008c360: 0f b7 10                    	movzwl	(%rax), %edx
14008c363: 83 ea 41                    	subl	$0x41, %edx
14008c366: 66 83 fa 19                 	cmpw	$0x19, %dx
14008c36a: 0f 87 e0 03 00 00           	ja	0x14008c750 <.text+0x8b750>
14008c370: 48 83 c0 02                 	addq	$0x2, %rax
14008c374: 48 39 c1                    	cmpq	%rax, %rcx
14008c377: 75 e7                       	jne	0x14008c360 <.text+0x8b360>
14008c379: 48 8d 54 24 60              	leaq	0x60(%rsp), %rdx
14008c37e: 48 8d 8c 24 40 02 00 00     	leaq	0x240(%rsp), %rcx
14008c386: e8 c5 ca ff ff              	callq	0x140088e50 <.text+0x87e50>
14008c38b: 4c 8b 84 24 48 02 00 00     	movq	0x248(%rsp), %r8
14008c393: 4c 8d 8c 24 30 02 00 00     	leaq	0x230(%rsp), %r9
14008c39b: 48 b8 7b 22 63 61 72 64 5f 6b       	movabsq	$0x6b5f64726163227b, %rax # imm = 0x6B5F64726163227B
14008c3a5: c6 84 24 3d 02 00 00 00     	movb	$0x0, 0x23d(%rsp)
14008c3ad: 48 89 84 24 30 02 00 00     	movq	%rax, 0x230(%rsp)
14008c3b5: 48 b8 64 5f 6b 65 79 22 3a 22       	movabsq	$0x223a2279656b5f64, %rax # imm = 0x223A2279656B5F64
14008c3bf: 48 8b 94 24 40 02 00 00     	movq	0x240(%rsp), %rdx
14008c3c7: 48 89 84 24 35 02 00 00     	movq	%rax, 0x235(%rsp)
14008c3cf: 49 8d 40 0d                 	leaq	0xd(%r8), %rax
14008c3d3: 4c 89 4c 24 38              	movq	%r9, 0x38(%rsp)
14008c3d8: 4c 89 8c 24 20 02 00 00     	movq	%r9, 0x220(%rsp)
14008c3e0: 48 c7 84 24 28 02 00 00 0d 00 00 00 	movq	$0xd, 0x228(%rsp)
14008c3ec: 48 83 f8 0f                 	cmpq	$0xf, %rax
14008c3f0: 76 2e                       	jbe	0x14008c420 <.text+0x8b420>
14008c3f2: 48 8d 8c 24 50 02 00 00     	leaq	0x250(%rsp), %rcx
14008c3fa: 48 39 ca                    	cmpq	%rcx, %rdx
14008c3fd: 74 21                       	je	0x14008c420 <.text+0x8b420>
14008c3ff: 48 39 84 24 50 02 00 00     	cmpq	%rax, 0x250(%rsp)
14008c407: 0f 83 3d 04 00 00           	jae	0x14008c84a <.text+0x8b84a>
14008c40d: 48 b8 f1 ff ff ff ff ff ff 7f       	movabsq	$0x7ffffffffffffff1, %rax # imm = 0x7FFFFFFFFFFFFFF1
14008c417: 4c 39 c0                    	cmpq	%r8, %rax
14008c41a: 0f 82 5d c1 04 00           	jb	0x1400d857d <.text+0xd757d>
14008c420: 48 8d 84 24 20 02 00 00     	leaq	0x220(%rsp), %rax
14008c428: 48 89 c1                    	movq	%rax, %rcx
14008c42b: 48 89 44 24 48              	movq	%rax, 0x48(%rsp)
14008c430: e8 9b 64 03 00              	callq	0x1400c28d0 <.text+0xc18d0>
14008c435: 4c 8d bc 24 70 02 00 00     	leaq	0x270(%rsp), %r15
14008c43d: 48 8d 50 10                 	leaq	0x10(%rax), %rdx
14008c441: 4c 89 bc 24 60 02 00 00     	movq	%r15, 0x260(%rsp)
14008c449: 48 8b 08                    	movq	(%rax), %rcx
14008c44c: 48 39 d1                    	cmpq	%rdx, %rcx
14008c44f: 0f 84 6b 04 00 00           	je	0x14008c8c0 <.text+0x8b8c0>
14008c455: 48 89 8c 24 60 02 00 00     	movq	%rcx, 0x260(%rsp)
14008c45d: 48 8b 48 10                 	movq	0x10(%rax), %rcx
14008c461: 48 89 8c 24 70 02 00 00     	movq	%rcx, 0x270(%rsp)
14008c469: 48 8b 48 08                 	movq	0x8(%rax), %rcx
14008c46d: 48 89 8c 24 68 02 00 00     	movq	%rcx, 0x268(%rsp)
14008c475: c6 40 10 00                 	movb	$0x0, 0x10(%rax)
14008c479: 48 c7 40 08 00 00 00 00     	movq	$0x0, 0x8(%rax)
14008c481: 48 89 10                    	movq	%rdx, (%rax)
14008c484: 48 b8 03 00 00 00 00 00 00 80       	movabsq	$-0x7ffffffffffffffd, %rax # imm = 0x8000000000000003
14008c48e: 48 03 84 24 68 02 00 00     	addq	0x268(%rsp), %rax
14008c496: 48 83 f8 01                 	cmpq	$0x1, %rax
14008c49a: 0f 86 17 c1 04 00           	jbe	0x1400d85b7 <.text+0xd75b7>
14008c4a0: 48 8d 84 24 60 02 00 00     	leaq	0x260(%rsp), %rax
14008c4a8: 41 b8 02 00 00 00           	movl	$0x2, %r8d
14008c4ae: 48 8d 15 60 e8 06 00        	leaq	0x6e860(%rip), %rdx     # 0x1400fad15
14008c4b5: 48 89 c1                    	movq	%rax, %rcx
14008c4b8: 48 89 44 24 30              	movq	%rax, 0x30(%rsp)
14008c4bd: e8 0e 64 03 00              	callq	0x1400c28d0 <.text+0xc18d0>
14008c4c2: 48 8d 94 24 90 00 00 00     	leaq	0x90(%rsp), %rdx
14008c4ca: 48 89 94 24 80 00 00 00     	movq	%rdx, 0x80(%rsp)
14008c4d2: 48 8b 08                    	movq	(%rax), %rcx
14008c4d5: 48 8d 50 10                 	leaq	0x10(%rax), %rdx
14008c4d9: 48 39 d1                    	cmpq	%rdx, %rcx
14008c4dc: 0f 84 0b 04 00 00           	je	0x14008c8ed <.text+0x8b8ed>
14008c4e2: 48 89 8c 24 80 00 00 00     	movq	%rcx, 0x80(%rsp)
14008c4ea: 48 8b 48 10                 	movq	0x10(%rax), %rcx
14008c4ee: 48 89 8c 24 90 00 00 00     	movq	%rcx, 0x90(%rsp)
14008c4f6: 48 8b 48 08                 	movq	0x8(%rax), %rcx
14008c4fa: 48 89 8c 24 88 00 00 00     	movq	%rcx, 0x88(%rsp)
14008c502: 48 89 10                    	movq	%rdx, (%rax)
14008c505: c6 40 10 00                 	movb	$0x0, 0x10(%rax)
14008c509: 48 8b 8c 24 60 02 00 00     	movq	0x260(%rsp), %rcx
14008c511: 48 c7 40 08 00 00 00 00     	movq	$0x0, 0x8(%rax)
14008c519: 4c 39 f9                    	cmpq	%r15, %rcx
14008c51c: 74 11                       	je	0x14008c52f <.text+0x8b52f>
14008c51e: 48 8b 84 24 70 02 00 00     	movq	0x270(%rsp), %rax
14008c526: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008c52a: e8 a1 db 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008c52f: 48 8b 8c 24 20 02 00 00     	movq	0x220(%rsp), %rcx
14008c537: 48 3b 4c 24 38              	cmpq	0x38(%rsp), %rcx
14008c53c: 74 11                       	je	0x14008c54f <.text+0x8b54f>
14008c53e: 48 8b 84 24 30 02 00 00     	movq	0x230(%rsp), %rax
14008c546: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008c54a: e8 81 db 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008c54f: 48 8b 8c 24 40 02 00 00     	movq	0x240(%rsp), %rcx
14008c557: 48 8d 84 24 50 02 00 00     	leaq	0x250(%rsp), %rax
14008c55f: 48 39 c1                    	cmpq	%rax, %rcx
14008c562: 74 11                       	je	0x14008c575 <.text+0x8b575>
14008c564: 48 8b 84 24 50 02 00 00     	movq	0x250(%rsp), %rax
14008c56c: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008c570: e8 5b db 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008c575: 48 8d 84 24 b0 00 00 00     	leaq	0xb0(%rsp), %rax
14008c57d: 45 31 d2                    	xorl	%r10d, %r10d
14008c580: 48 8d 0d a1 e8 06 00        	leaq	0x6e8a1(%rip), %rcx     # 0x1400fae28
14008c587: c6 84 24 b0 00 00 00 00     	movb	$0x0, 0xb0(%rsp)
14008c58f: 48 89 84 24 a0 00 00 00     	movq	%rax, 0xa0(%rsp)
14008c597: 48 8d 84 24 d0 00 00 00     	leaq	0xd0(%rsp), %rax
14008c59f: 48 c7 84 24 a8 00 00 00 00 00 00 00 	movq	$0x0, 0xa8(%rsp)
14008c5ab: 48 89 84 24 c0 00 00 00     	movq	%rax, 0xc0(%rsp)
14008c5b3: 66 44 89 94 24 d0 00 00 00  	movw	%r10w, 0xd0(%rsp)
14008c5bc: 48 c7 84 24 c8 00 00 00 00 00 00 00 	movq	$0x0, 0xc8(%rsp)
14008c5c8: 4c 89 bc 24 60 02 00 00     	movq	%r15, 0x260(%rsp)
14008c5d0: e8 ab 90 01 00              	callq	0x1400a5680 <.text+0xa4680>
14008c5d5: 4c 8d 0c 00                 	leaq	(%rax,%rax), %r9
14008c5d9: 4c 89 ca                    	movq	%r9, %rdx
14008c5dc: 48 d1 fa                    	sarq	%rdx
14008c5df: 48 89 94 24 40 02 00 00     	movq	%rdx, 0x240(%rsp)
14008c5e7: 49 83 f9 0e                 	cmpq	$0xe, %r9
14008c5eb: 0f 87 8b 01 00 00           	ja	0x14008c77c <.text+0x8b77c>
14008c5f1: 48 8b 8c 24 60 02 00 00     	movq	0x260(%rsp), %rcx
14008c5f9: 48 83 fa 01                 	cmpq	$0x1, %rdx
14008c5fd: 0f 84 33 02 00 00           	je	0x14008c836 <.text+0x8b836>
14008c603: 31 c0                       	xorl	%eax, %eax
14008c605: 48 85 d2                    	testq	%rdx, %rdx
14008c608: 0f 85 cf 01 00 00           	jne	0x14008c7dd <.text+0x8b7dd>
14008c60e: 45 31 c0                    	xorl	%r8d, %r8d
14008c611: 66 44 89 04 01              	movw	%r8w, (%rcx,%rax)
14008c616: 48 8d 8c 24 40 02 00 00     	leaq	0x240(%rsp), %rcx
14008c61e: 48 89 94 24 68 02 00 00     	movq	%rdx, 0x268(%rsp)
14008c626: e8 e5 d6 ff ff              	callq	0x140089d10 <.text+0x88d10>
14008c62b: 48 8d 84 24 c0 00 00 00     	leaq	0xc0(%rsp), %rax
14008c633: 48 8b 54 24 30              	movq	0x30(%rsp), %rdx
14008c638: 4c 8d 8c 24 a0 00 00 00     	leaq	0xa0(%rsp), %r9
14008c640: 48 89 44 24 20              	movq	%rax, 0x20(%rsp)
14008c645: 4c 8d 84 24 80 00 00 00     	leaq	0x80(%rsp), %r8
14008c64d: 48 8d 8c 24 40 02 00 00     	leaq	0x240(%rsp), %rcx
14008c655: 4c 89 4c 24 40              	movq	%r9, 0x40(%rsp)
14008c65a: 4c 89 44 24 58              	movq	%r8, 0x58(%rsp)
14008c65f: 48 89 44 24 50              	movq	%rax, 0x50(%rsp)
14008c664: e8 87 d0 ff ff              	callq	0x1400896f0 <.text+0x886f0>
14008c669: 48 8b 8c 24 40 02 00 00     	movq	0x240(%rsp), %rcx
14008c671: 88 44 24 48                 	movb	%al, 0x48(%rsp)
14008c675: 48 8d 84 24 50 02 00 00     	leaq	0x250(%rsp), %rax
14008c67d: 48 39 c1                    	cmpq	%rax, %rcx
14008c680: 74 12                       	je	0x14008c694 <.text+0x8b694>
14008c682: 48 8b 84 24 50 02 00 00     	movq	0x250(%rsp), %rax
14008c68a: 48 8d 54 00 02              	leaq	0x2(%rax,%rax), %rdx
14008c68f: e8 3c da 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008c694: 48 8b 8c 24 60 02 00 00     	movq	0x260(%rsp), %rcx
14008c69c: 4c 39 f9                    	cmpq	%r15, %rcx
14008c69f: 74 12                       	je	0x14008c6b3 <.text+0x8b6b3>
14008c6a1: 48 8b 84 24 70 02 00 00     	movq	0x270(%rsp), %rax
14008c6a9: 48 8d 54 00 02              	leaq	0x2(%rax,%rax), %rdx
14008c6ae: e8 1d da 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008c6b3: 80 7c 24 48 00              	cmpb	$0x0, 0x48(%rsp)
14008c6b8: 0f 85 4a 01 00 00           	jne	0x14008c808 <.text+0x8b808>
14008c6be: 48 8b 54 24 50              	movq	0x50(%rsp), %rdx
14008c6c3: 48 8d 4b 08                 	leaq	0x8(%rbx), %rcx
14008c6c7: e8 f4 9f 03 00              	callq	0x1400c66c0 <.text+0xc56c0>
14008c6cc: 48 8b 8c 24 c0 00 00 00     	movq	0xc0(%rsp), %rcx
14008c6d4: 48 8d 84 24 d0 00 00 00     	leaq	0xd0(%rsp), %rax
14008c6dc: 48 39 c1                    	cmpq	%rax, %rcx
14008c6df: 74 12                       	je	0x14008c6f3 <.text+0x8b6f3>
14008c6e1: 48 8b 84 24 d0 00 00 00     	movq	0xd0(%rsp), %rax
14008c6e9: 48 8d 54 00 02              	leaq	0x2(%rax,%rax), %rdx
14008c6ee: e8 dd d9 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008c6f3: 48 8b 8c 24 a0 00 00 00     	movq	0xa0(%rsp), %rcx
14008c6fb: 48 8d 84 24 b0 00 00 00     	leaq	0xb0(%rsp), %rax
14008c703: 48 39 c1                    	cmpq	%rax, %rcx
14008c706: 74 11                       	je	0x14008c719 <.text+0x8b719>
14008c708: 48 8b 84 24 b0 00 00 00     	movq	0xb0(%rsp), %rax
14008c710: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008c714: e8 b7 d9 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008c719: 48 8b 8c 24 80 00 00 00     	movq	0x80(%rsp), %rcx
14008c721: 48 8d 84 24 90 00 00 00     	leaq	0x90(%rsp), %rax
14008c729: 48 39 c1                    	cmpq	%rax, %rcx
14008c72c: 0f 84 98 fb ff ff           	je	0x14008c2ca <.text+0x8b2ca>
14008c732: 48 8b 84 24 90 00 00 00     	movq	0x90(%rsp), %rax
14008c73a: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008c73e: e8 8d d9 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008c743: e9 82 fb ff ff              	jmp	0x14008c2ca <.text+0x8b2ca>
14008c748: 0f 1f 84 00 00 00 00 00     	nopl	(%rax,%rax)
14008c750: 48 8d 0d 0f e6 06 00        	leaq	0x6e60f(%rip), %rcx     # 0x1400fad66
14008c757: e8 24 8f 01 00              	callq	0x1400a5680 <.text+0xa4680>
14008c75c: 48 8d 4b 08                 	leaq	0x8(%rbx), %rcx
14008c760: 4c 8d 0d ff e5 06 00        	leaq	0x6e5ff(%rip), %r9      # 0x1400fad66
14008c767: 31 d2                       	xorl	%edx, %edx
14008c769: 48 89 44 24 20              	movq	%rax, 0x20(%rsp)
14008c76e: 4c 8b 43 10                 	movq	0x10(%rbx), %r8
14008c772: e8 a9 75 03 00              	callq	0x1400c3d20 <.text+0xc2d20>
14008c777: e9 4e fb ff ff              	jmp	0x14008c2ca <.text+0x8b2ca>
14008c77c: 48 8d 84 24 a0 00 00 00     	leaq	0xa0(%rsp), %rax
14008c784: 48 8b 4c 24 30              	movq	0x30(%rsp), %rcx
14008c789: 45 31 c0                    	xorl	%r8d, %r8d
14008c78c: 48 8d 94 24 40 02 00 00     	leaq	0x240(%rsp), %rdx
14008c794: 48 89 44 24 40              	movq	%rax, 0x40(%rsp)
14008c799: 48 8d 84 24 80 00 00 00     	leaq	0x80(%rsp), %rax
14008c7a1: 48 89 44 24 58              	movq	%rax, 0x58(%rsp)
14008c7a6: 48 8d 84 24 c0 00 00 00     	leaq	0xc0(%rsp), %rax
14008c7ae: 4c 89 4c 24 48              	movq	%r9, 0x48(%rsp)
14008c7b3: 48 89 44 24 50              	movq	%rax, 0x50(%rsp)
14008c7b8: e8 23 a0 03 00              	callq	0x1400c67e0 <.text+0xc57e0>
14008c7bd: 48 89 84 24 60 02 00 00     	movq	%rax, 0x260(%rsp)
14008c7c5: 48 89 c1                    	movq	%rax, %rcx
14008c7c8: 4c 8b 4c 24 48              	movq	0x48(%rsp), %r9
14008c7cd: 48 8b 84 24 40 02 00 00     	movq	0x240(%rsp), %rax
14008c7d5: 48 89 84 24 70 02 00 00     	movq	%rax, 0x270(%rsp)
14008c7dd: 48 8d 15 44 e6 06 00        	leaq	0x6e644(%rip), %rdx     # 0x1400fae28
14008c7e4: 4d 89 c8                    	movq	%r9, %r8
14008c7e7: e8 fc 8f 01 00              	callq	0x1400a57e8 <.text+0xa47e8>
14008c7ec: 48 8b 94 24 40 02 00 00     	movq	0x240(%rsp), %rdx
14008c7f4: 48 8b 8c 24 60 02 00 00     	movq	0x260(%rsp), %rcx
14008c7fc: 48 8d 04 12                 	leaq	(%rdx,%rdx), %rax
14008c800: e9 09 fe ff ff              	jmp	0x14008c60e <.text+0x8b60e>
14008c805: 0f 1f 00                    	nopl	(%rax)
14008c808: 48 8b 4c 24 40              	movq	0x40(%rsp), %rcx
14008c80d: e8 6e ce ff ff              	callq	0x140089680 <.text+0x88680>
14008c812: 84 c0                       	testb	%al, %al
14008c814: 0f 84 04 01 00 00           	je	0x14008c91e <.text+0x8b91e>
14008c81a: 48 8b 43 08                 	movq	0x8(%rbx), %rax
14008c81e: 45 31 f6                    	xorl	%r14d, %r14d
14008c821: c6 43 01 01                 	movb	$0x1, 0x1(%rbx)
14008c825: 66 44 89 30                 	movw	%r14w, (%rax)
14008c829: 48 c7 43 10 00 00 00 00     	movq	$0x0, 0x10(%rbx)
14008c831: e9 96 fe ff ff              	jmp	0x14008c6cc <.text+0x8b6cc>
14008c836: 41 b9 2f 00 00 00           	movl	$0x2f, %r9d
14008c83c: b8 02 00 00 00              	movl	$0x2, %eax
14008c841: 66 44 89 09                 	movw	%r9w, (%rcx)
14008c845: e9 c4 fd ff ff              	jmp	0x14008c60e <.text+0x8b60e>
14008c84a: 48 c7 44 24 20 0d 00 00 00  	movq	$0xd, 0x20(%rsp)
14008c853: 48 8d 84 24 20 02 00 00     	leaq	0x220(%rsp), %rax
14008c85b: 45 31 c0                    	xorl	%r8d, %r8d
14008c85e: 31 d2                       	xorl	%edx, %edx
14008c860: 48 8d 8c 24 40 02 00 00     	leaq	0x240(%rsp), %rcx
14008c868: 48 89 44 24 48              	movq	%rax, 0x48(%rsp)
14008c86d: e8 5e 30 03 00              	callq	0x1400bf8d0 <.text+0xbe8d0>
14008c872: 4c 8d bc 24 70 02 00 00     	leaq	0x270(%rsp), %r15
14008c87a: 48 8d 50 10                 	leaq	0x10(%rax), %rdx
14008c87e: 4c 89 bc 24 60 02 00 00     	movq	%r15, 0x260(%rsp)
14008c886: 48 8b 08                    	movq	(%rax), %rcx
14008c889: 48 39 d1                    	cmpq	%rdx, %rcx
14008c88c: 0f 85 c3 fb ff ff           	jne	0x14008c455 <.text+0x8b455>
14008c892: 48 8b 48 08                 	movq	0x8(%rax), %rcx
14008c896: 4c 8d 41 01                 	leaq	0x1(%rcx), %r8
14008c89a: 41 83 e0 3f                 	andl	$0x3f, %r8d
14008c89e: 0f 84 c9 fb ff ff           	je	0x14008c46d <.text+0x8b46d>
14008c8a4: 31 c9                       	xorl	%ecx, %ecx
14008c8a6: 41 89 c9                    	movl	%ecx, %r9d
14008c8a9: 83 c1 01                    	addl	$0x1, %ecx
14008c8ac: 46 0f b6 14 0a              	movzbl	(%rdx,%r9), %r10d
14008c8b1: 47 88 14 0f                 	movb	%r10b, (%r15,%r9)
14008c8b5: 44 39 c1                    	cmpl	%r8d, %ecx
14008c8b8: 72 ec                       	jb	0x14008c8a6 <.text+0x8b8a6>
14008c8ba: e9 aa fb ff ff              	jmp	0x14008c469 <.text+0x8b469>
14008c8bf: 90                          	nop
14008c8c0: 48 8b 48 08                 	movq	0x8(%rax), %rcx
14008c8c4: 4c 8d 41 01                 	leaq	0x1(%rcx), %r8
14008c8c8: 41 83 e0 3f                 	andl	$0x3f, %r8d
14008c8cc: 0f 84 9b fb ff ff           	je	0x14008c46d <.text+0x8b46d>
14008c8d2: 31 c9                       	xorl	%ecx, %ecx
14008c8d4: 41 89 c9                    	movl	%ecx, %r9d
14008c8d7: 83 c1 01                    	addl	$0x1, %ecx
14008c8da: 46 0f b6 14 0a              	movzbl	(%rdx,%r9), %r10d
14008c8df: 47 88 14 0f                 	movb	%r10b, (%r15,%r9)
14008c8e3: 44 39 c1                    	cmpl	%r8d, %ecx
14008c8e6: 72 ec                       	jb	0x14008c8d4 <.text+0x8b8d4>
14008c8e8: e9 7c fb ff ff              	jmp	0x14008c469 <.text+0x8b469>
14008c8ed: 48 8b 48 08                 	movq	0x8(%rax), %rcx
14008c8f1: 4c 8d 41 01                 	leaq	0x1(%rcx), %r8
14008c8f5: 41 83 e0 3f                 	andl	$0x3f, %r8d
14008c8f9: 0f 84 fb fb ff ff           	je	0x14008c4fa <.text+0x8b4fa>
14008c8ff: 31 c9                       	xorl	%ecx, %ecx
14008c901: 41 89 c9                    	movl	%ecx, %r9d
14008c904: 83 c1 01                    	addl	$0x1, %ecx
14008c907: 46 0f b6 14 0a              	movzbl	(%rdx,%r9), %r10d
14008c90c: 46 88 94 0c 90 00 00 00     	movb	%r10b, 0x90(%rsp,%r9)
14008c914: 44 39 c1                    	cmpl	%r8d, %ecx
14008c917: 72 e8                       	jb	0x14008c901 <.text+0x8b901>
14008c919: e9 d8 fb ff ff              	jmp	0x14008c4f6 <.text+0x8b4f6>
14008c91e: 48 8b 54 24 40              	movq	0x40(%rsp), %rdx
14008c923: 4c 8d 05 70 e4 06 00        	leaq	0x6e470(%rip), %r8      # 0x1400fad9a
14008c92a: 48 8d 8c 24 e0 00 00 00     	leaq	0xe0(%rsp), %rcx
14008c932: e8 09 ca ff ff              	callq	0x140089340 <.text+0x88340>
14008c937: 48 83 bc 24 e8 00 00 00 04  	cmpq	$0x4, 0xe8(%rsp)
14008c940: 0f 84 36 01 00 00           	je	0x14008ca7c <.text+0x8ba7c>
14008c946: 48 8b 54 24 40              	movq	0x40(%rsp), %rdx
14008c94b: 4c 8d 05 4b e4 06 00        	leaq	0x6e44b(%rip), %r8      # 0x1400fad9d
14008c952: 48 8d 8c 24 40 02 00 00     	leaq	0x240(%rsp), %rcx
14008c95a: e8 e1 c9 ff ff              	callq	0x140089340 <.text+0x88340>
14008c95f: 48 83 bc 24 48 02 00 00 00  	cmpq	$0x0, 0x248(%rsp)
14008c968: 0f 85 dc 04 00 00           	jne	0x14008ce4a <.text+0x8be4a>
14008c96e: 48 8d 0d d9 e4 06 00        	leaq	0x6e4d9(%rip), %rcx     # 0x1400fae4e
14008c975: 4c 89 bc 24 60 02 00 00     	movq	%r15, 0x260(%rsp)
14008c97d: e8 fe 8c 01 00              	callq	0x1400a5680 <.text+0xa4680>
14008c982: 48 8d 34 00                 	leaq	(%rax,%rax), %rsi
14008c986: 49 89 f0                    	movq	%rsi, %r8
14008c989: 49 d1 f8                    	sarq	%r8
14008c98c: 4c 89 84 24 20 02 00 00     	movq	%r8, 0x220(%rsp)
14008c994: 48 83 fe 0e                 	cmpq	$0xe, %rsi
14008c998: 0f 87 15 08 00 00           	ja	0x14008d1b3 <.text+0x8c1b3>
14008c99e: 48 8b 8c 24 60 02 00 00     	movq	0x260(%rsp), %rcx
14008c9a6: 49 83 f8 01                 	cmpq	$0x1, %r8
14008c9aa: 0f 84 ef 07 00 00           	je	0x14008d19f <.text+0x8c19f>
14008c9b0: 31 c0                       	xorl	%eax, %eax
14008c9b2: 4d 85 c0                    	testq	%r8, %r8
14008c9b5: 0f 85 28 08 00 00           	jne	0x14008d1e3 <.text+0x8c1e3>
14008c9bb: 45 31 e4                    	xorl	%r12d, %r12d
14008c9be: 66 44 89 24 01              	movw	%r12w, (%rcx,%rax)
14008c9c3: 4c 89 84 24 68 02 00 00     	movq	%r8, 0x268(%rsp)
14008c9cb: 48 8b 43 08                 	movq	0x8(%rbx), %rax
14008c9cf: 48 39 c7                    	cmpq	%rax, %rdi
14008c9d2: 0f 84 a9 07 00 00           	je	0x14008d181 <.text+0x8c181>
14008c9d8: 4c 39 f9                    	cmpq	%r15, %rcx
14008c9db: 0f 84 5b 07 00 00           	je	0x14008d13c <.text+0x8c13c>
14008c9e1: 48 89 4b 08                 	movq	%rcx, 0x8(%rbx)
14008c9e5: 48 8b 8c 24 70 02 00 00     	movq	0x270(%rsp), %rcx
14008c9ed: 48 8b 53 18                 	movq	0x18(%rbx), %rdx
14008c9f1: 4c 89 43 10                 	movq	%r8, 0x10(%rbx)
14008c9f5: 48 89 4b 18                 	movq	%rcx, 0x18(%rbx)
14008c9f9: 48 85 c0                    	testq	%rax, %rax
14008c9fc: 74 29                       	je	0x14008ca27 <.text+0x8ba27>
14008c9fe: 48 89 94 24 70 02 00 00     	movq	%rdx, 0x270(%rsp)
14008ca06: 45 31 db                    	xorl	%r11d, %r11d
14008ca09: 66 44 89 18                 	movw	%r11w, (%rax)
14008ca0d: 4c 39 f8                    	cmpq	%r15, %rax
14008ca10: 74 15                       	je	0x14008ca27 <.text+0x8ba27>
14008ca12: 48 8b 94 24 70 02 00 00     	movq	0x270(%rsp), %rdx
14008ca1a: 48 89 c1                    	movq	%rax, %rcx
14008ca1d: 48 8d 54 12 02              	leaq	0x2(%rdx,%rdx), %rdx
14008ca22: e8 a9 d6 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008ca27: 48 8b 8c 24 40 02 00 00     	movq	0x240(%rsp), %rcx
14008ca2f: 48 8d 84 24 50 02 00 00     	leaq	0x250(%rsp), %rax
14008ca37: 48 39 c1                    	cmpq	%rax, %rcx
14008ca3a: 74 11                       	je	0x14008ca4d <.text+0x8ba4d>
14008ca3c: 48 8b 84 24 50 02 00 00     	movq	0x250(%rsp), %rax
14008ca44: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008ca48: e8 83 d6 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008ca4d: 48 8b 8c 24 e0 00 00 00     	movq	0xe0(%rsp), %rcx
14008ca55: 48 8d 84 24 f0 00 00 00     	leaq	0xf0(%rsp), %rax
14008ca5d: 48 39 c1                    	cmpq	%rax, %rcx
14008ca60: 0f 84 66 fc ff ff           	je	0x14008c6cc <.text+0x8b6cc>
14008ca66: 48 8b 84 24 f0 00 00 00     	movq	0xf0(%rsp), %rax
14008ca6e: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008ca72: e8 59 d6 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008ca77: e9 50 fc ff ff              	jmp	0x14008c6cc <.text+0x8b6cc>
14008ca7c: 48 8b 84 24 e0 00 00 00     	movq	0xe0(%rsp), %rax
14008ca84: 81 38 74 72 75 65           	cmpl	$0x65757274, (%rax)     # imm = 0x65757274
14008ca8a: 0f 85 b6 fe ff ff           	jne	0x14008c946 <.text+0x8b946>
14008ca90: 48 8b 54 24 40              	movq	0x40(%rsp), %rdx
14008ca95: 4c 8d 05 ae e3 06 00        	leaq	0x6e3ae(%rip), %r8      # 0x1400fae4a
14008ca9c: 48 8d 8c 24 00 01 00 00     	leaq	0x100(%rsp), %rcx
14008caa4: e8 97 c8 ff ff              	callq	0x140089340 <.text+0x88340>
14008caa9: 48 8b 7c 24 40              	movq	0x40(%rsp), %rdi
14008caae: 4c 8d 05 a3 e3 06 00        	leaq	0x6e3a3(%rip), %r8      # 0x1400fae58
14008cab5: 48 8d 8c 24 20 01 00 00     	leaq	0x120(%rsp), %rcx
14008cabd: 48 89 fa                    	movq	%rdi, %rdx
14008cac0: e8 7b c8 ff ff              	callq	0x140089340 <.text+0x88340>
14008cac5: 4c 8d 05 8f e3 06 00        	leaq	0x6e38f(%rip), %r8      # 0x1400fae5b
14008cacc: 48 89 fa                    	movq	%rdi, %rdx
14008cacf: 48 8d 8c 24 40 01 00 00     	leaq	0x140(%rsp), %rcx
14008cad7: e8 64 c8 ff ff              	callq	0x140089340 <.text+0x88340>
14008cadc: 4c 8d 05 7c e3 06 00        	leaq	0x6e37c(%rip), %r8      # 0x1400fae5f
14008cae3: 48 89 fa                    	movq	%rdi, %rdx
14008cae6: 48 8d 8c 24 60 01 00 00     	leaq	0x160(%rsp), %rcx
14008caee: e8 4d c8 ff ff              	callq	0x140089340 <.text+0x88340>
14008caf3: 4c 8d 05 b6 e2 06 00        	leaq	0x6e2b6(%rip), %r8      # 0x1400fadb0
14008cafa: 48 89 fa                    	movq	%rdi, %rdx
14008cafd: 48 8d 8c 24 80 01 00 00     	leaq	0x180(%rsp), %rcx
14008cb05: e8 36 c8 ff ff              	callq	0x140089340 <.text+0x88340>
14008cb0a: 4c 8d 05 a2 e2 06 00        	leaq	0x6e2a2(%rip), %r8      # 0x1400fadb3
14008cb11: 48 89 fa                    	movq	%rdi, %rdx
14008cb14: 48 8d 8c 24 a0 01 00 00     	leaq	0x1a0(%rsp), %rcx
14008cb1c: e8 1f c8 ff ff              	callq	0x140089340 <.text+0x88340>
14008cb21: 4c 8d 05 ac e2 06 00        	leaq	0x6e2ac(%rip), %r8      # 0x1400fadd4
14008cb28: 48 89 fa                    	movq	%rdi, %rdx
14008cb2b: 48 8d 8c 24 c0 01 00 00     	leaq	0x1c0(%rsp), %rcx
14008cb33: e8 08 c8 ff ff              	callq	0x140089340 <.text+0x88340>
14008cb38: 4c 8d 05 9a e2 06 00        	leaq	0x6e29a(%rip), %r8      # 0x1400fadd9
14008cb3f: 48 89 fa                    	movq	%rdi, %rdx
14008cb42: 48 8d 8c 24 e0 01 00 00     	leaq	0x1e0(%rsp), %rcx
14008cb4a: e8 f1 c7 ff ff              	callq	0x140089340 <.text+0x88340>
14008cb4f: 4c 8d 05 0e e3 06 00        	leaq	0x6e30e(%rip), %r8      # 0x1400fae64
14008cb56: 48 89 fa                    	movq	%rdi, %rdx
14008cb59: 48 8d 8c 24 00 02 00 00     	leaq	0x200(%rsp), %rcx
14008cb61: e8 da c7 ff ff              	callq	0x140089340 <.text+0x88340>
14008cb66: 48 8b 7c 24 38              	movq	0x38(%rsp), %rdi
14008cb6b: 48 8d 84 24 20 02 00 00     	leaq	0x220(%rsp), %rax
14008cb73: 45 31 d2                    	xorl	%r10d, %r10d
14008cb76: 4c 8d 8c 24 60 01 00 00     	leaq	0x160(%rsp), %r9
14008cb7e: 48 89 44 24 20              	movq	%rax, 0x20(%rsp)
14008cb83: 4c 8d 84 24 40 01 00 00     	leaq	0x140(%rsp), %r8
14008cb8b: 48 8d 94 24 20 01 00 00     	leaq	0x120(%rsp), %rdx
14008cb93: 48 89 bc 24 20 02 00 00     	movq	%rdi, 0x220(%rsp)
14008cb9b: 48 8d bc 24 50 02 00 00     	leaq	0x250(%rsp), %rdi
14008cba3: 48 8d 8c 24 00 01 00 00     	leaq	0x100(%rsp), %rcx
14008cbab: 48 89 bc 24 40 02 00 00     	movq	%rdi, 0x240(%rsp)
14008cbb3: 48 8d bc 24 40 02 00 00     	leaq	0x240(%rsp), %rdi
14008cbbb: 48 89 7c 24 28              	movq	%rdi, 0x28(%rsp)
14008cbc0: 48 89 44 24 48              	movq	%rax, 0x48(%rsp)
14008cbc5: c6 84 24 30 02 00 00 00     	movb	$0x0, 0x230(%rsp)
14008cbcd: 48 c7 84 24 28 02 00 00 00 00 00 00 	movq	$0x0, 0x228(%rsp)
14008cbd9: 66 44 89 94 24 50 02 00 00  	movw	%r10w, 0x250(%rsp)
14008cbe2: 48 c7 84 24 48 02 00 00 00 00 00 00 	movq	$0x0, 0x248(%rsp)
14008cbee: e8 0d 0b 00 00              	callq	0x14008d700 <.text+0x8c700>
14008cbf3: 84 c0                       	testb	%al, %al
14008cbf5: 0f 85 b2 01 00 00           	jne	0x14008cdad <.text+0x8bdad>
14008cbfb: 48 8d 4b 08                 	leaq	0x8(%rbx), %rcx
14008cbff: 48 89 fa                    	movq	%rdi, %rdx
14008cc02: e8 b9 9a 03 00              	callq	0x1400c66c0 <.text+0xc56c0>
14008cc07: 48 8b 8c 24 40 02 00 00     	movq	0x240(%rsp), %rcx
14008cc0f: 48 8d 84 24 50 02 00 00     	leaq	0x250(%rsp), %rax
14008cc17: 48 39 c1                    	cmpq	%rax, %rcx
14008cc1a: 74 12                       	je	0x14008cc2e <.text+0x8bc2e>
14008cc1c: 48 8b 84 24 50 02 00 00     	movq	0x250(%rsp), %rax
14008cc24: 48 8d 54 00 02              	leaq	0x2(%rax,%rax), %rdx
14008cc29: e8 a2 d4 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008cc2e: 48 8b 8c 24 20 02 00 00     	movq	0x220(%rsp), %rcx
14008cc36: 48 3b 4c 24 38              	cmpq	0x38(%rsp), %rcx
14008cc3b: 74 11                       	je	0x14008cc4e <.text+0x8bc4e>
14008cc3d: 48 8b 84 24 30 02 00 00     	movq	0x230(%rsp), %rax
14008cc45: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008cc49: e8 82 d4 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008cc4e: 48 8b 8c 24 00 02 00 00     	movq	0x200(%rsp), %rcx
14008cc56: 48 8d 84 24 10 02 00 00     	leaq	0x210(%rsp), %rax
14008cc5e: 48 39 c1                    	cmpq	%rax, %rcx
14008cc61: 74 11                       	je	0x14008cc74 <.text+0x8bc74>
14008cc63: 48 8b 84 24 10 02 00 00     	movq	0x210(%rsp), %rax
14008cc6b: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008cc6f: e8 5c d4 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008cc74: 48 8b 8c 24 e0 01 00 00     	movq	0x1e0(%rsp), %rcx
14008cc7c: 48 8d 84 24 f0 01 00 00     	leaq	0x1f0(%rsp), %rax
14008cc84: 48 39 c1                    	cmpq	%rax, %rcx
14008cc87: 74 11                       	je	0x14008cc9a <.text+0x8bc9a>
14008cc89: 48 8b 84 24 f0 01 00 00     	movq	0x1f0(%rsp), %rax
14008cc91: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008cc95: e8 36 d4 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008cc9a: 48 8b 8c 24 c0 01 00 00     	movq	0x1c0(%rsp), %rcx
14008cca2: 48 8d 84 24 d0 01 00 00     	leaq	0x1d0(%rsp), %rax
14008ccaa: 48 39 c1                    	cmpq	%rax, %rcx
14008ccad: 74 11                       	je	0x14008ccc0 <.text+0x8bcc0>
14008ccaf: 48 8b 84 24 d0 01 00 00     	movq	0x1d0(%rsp), %rax
14008ccb7: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008ccbb: e8 10 d4 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008ccc0: 48 8b 8c 24 a0 01 00 00     	movq	0x1a0(%rsp), %rcx
14008ccc8: 48 8d 84 24 b0 01 00 00     	leaq	0x1b0(%rsp), %rax
14008ccd0: 48 39 c1                    	cmpq	%rax, %rcx
14008ccd3: 74 11                       	je	0x14008cce6 <.text+0x8bce6>
14008ccd5: 48 8b 84 24 b0 01 00 00     	movq	0x1b0(%rsp), %rax
14008ccdd: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008cce1: e8 ea d3 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008cce6: 48 8b 8c 24 80 01 00 00     	movq	0x180(%rsp), %rcx
14008ccee: 48 8d 84 24 90 01 00 00     	leaq	0x190(%rsp), %rax
14008ccf6: 48 39 c1                    	cmpq	%rax, %rcx
14008ccf9: 74 11                       	je	0x14008cd0c <.text+0x8bd0c>
14008ccfb: 48 8b 84 24 90 01 00 00     	movq	0x190(%rsp), %rax
14008cd03: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008cd07: e8 c4 d3 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008cd0c: 48 8b 8c 24 60 01 00 00     	movq	0x160(%rsp), %rcx
14008cd14: 48 8d 84 24 70 01 00 00     	leaq	0x170(%rsp), %rax
14008cd1c: 48 39 c1                    	cmpq	%rax, %rcx
14008cd1f: 74 11                       	je	0x14008cd32 <.text+0x8bd32>
14008cd21: 48 8b 84 24 70 01 00 00     	movq	0x170(%rsp), %rax
14008cd29: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008cd2d: e8 9e d3 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008cd32: 48 8b 8c 24 40 01 00 00     	movq	0x140(%rsp), %rcx
14008cd3a: 48 8d 84 24 50 01 00 00     	leaq	0x150(%rsp), %rax
14008cd42: 48 39 c1                    	cmpq	%rax, %rcx
14008cd45: 74 11                       	je	0x14008cd58 <.text+0x8bd58>
14008cd47: 48 8b 84 24 50 01 00 00     	movq	0x150(%rsp), %rax
14008cd4f: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008cd53: e8 78 d3 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008cd58: 48 8b 8c 24 20 01 00 00     	movq	0x120(%rsp), %rcx
14008cd60: 48 8d 84 24 30 01 00 00     	leaq	0x130(%rsp), %rax
14008cd68: 48 39 c1                    	cmpq	%rax, %rcx
14008cd6b: 74 11                       	je	0x14008cd7e <.text+0x8bd7e>
14008cd6d: 48 8b 84 24 30 01 00 00     	movq	0x130(%rsp), %rax
14008cd75: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008cd79: e8 52 d3 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008cd7e: 48 8b 8c 24 00 01 00 00     	movq	0x100(%rsp), %rcx
14008cd86: 48 8d 84 24 10 01 00 00     	leaq	0x110(%rsp), %rax
14008cd8e: 48 39 c1                    	cmpq	%rax, %rcx
14008cd91: 0f 84 b6 fc ff ff           	je	0x14008ca4d <.text+0x8ba4d>
14008cd97: 48 8b 84 24 10 01 00 00     	movq	0x110(%rsp), %rax
14008cd9f: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008cda3: e8 28 d3 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008cda8: e9 a0 fc ff ff              	jmp	0x14008ca4d <.text+0x8ba4d>
14008cdad: 48 8b 4c 24 48              	movq	0x48(%rsp), %rcx
14008cdb2: 45 31 c0                    	xorl	%r8d, %r8d
14008cdb5: 41 b9 0a 00 00 00           	movl	$0xa, %r9d
14008cdbb: 48 8d 15 a9 e0 06 00        	leaq	0x6e0a9(%rip), %rdx     # 0x1400fae6b
14008cdc2: e8 19 a6 02 00              	callq	0x1400b73e0 <.text+0xb63e0>
14008cdc7: 48 8b 4c 24 48              	movq	0x48(%rsp), %rcx
14008cdcc: 45 31 c0                    	xorl	%r8d, %r8d
14008cdcf: 41 b9 08 00 00 00           	movl	$0x8, %r9d
14008cdd5: 48 8d 15 9a e0 06 00        	leaq	0x6e09a(%rip), %rdx     # 0x1400fae76
14008cddc: 48 89 c7                    	movq	%rax, %rdi
14008cddf: e8 fc a5 02 00              	callq	0x1400b73e0 <.text+0xb63e0>
14008cde4: 48 83 f8 ff                 	cmpq	$-0x1, %rax
14008cde8: 0f 84 a7 00 00 00           	je	0x14008ce95 <.text+0x8be95>
14008cdee: 48 8b 4c 24 48              	movq	0x48(%rsp), %rcx
14008cdf3: 45 31 c0                    	xorl	%r8d, %r8d
14008cdf6: 41 b9 05 00 00 00           	movl	$0x5, %r9d
14008cdfc: 48 8d 15 7c e0 06 00        	leaq	0x6e07c(%rip), %rdx     # 0x1400fae7f
14008ce03: e8 d8 a5 02 00              	callq	0x1400b73e0 <.text+0xb63e0>
14008ce08: 48 83 f8 ff                 	cmpq	$-0x1, %rax
14008ce0c: 74 63                       	je	0x14008ce71 <.text+0x8be71>
14008ce0e: 48 83 ff ff                 	cmpq	$-0x1, %rdi
14008ce12: 0f 85 8e 00 00 00           	jne	0x14008cea6 <.text+0x8bea6>
14008ce18: 4c 8d 0d c9 df 06 00        	leaq	0x6dfc9(%rip), %r9      # 0x1400fade8
14008ce1f: 4c 89 c9                    	movq	%r9, %rcx
14008ce22: 4c 89 4c 24 30              	movq	%r9, 0x30(%rsp)
14008ce27: e8 54 88 01 00              	callq	0x1400a5680 <.text+0xa4680>
14008ce2c: 4c 8b 4c 24 30              	movq	0x30(%rsp), %r9
14008ce31: 48 8d 4b 08                 	leaq	0x8(%rbx), %rcx
14008ce35: 31 d2                       	xorl	%edx, %edx
14008ce37: 48 89 44 24 20              	movq	%rax, 0x20(%rsp)
14008ce3c: 4c 8b 43 10                 	movq	0x10(%rbx), %r8
14008ce40: e8 db 6e 03 00              	callq	0x1400c3d20 <.text+0xc2d20>
14008ce45: e9 bd fd ff ff              	jmp	0x14008cc07 <.text+0x8bc07>
14008ce4a: 48 8b 4c 24 30              	movq	0x30(%rsp), %rcx
14008ce4f: 48 8d 94 24 40 02 00 00     	leaq	0x240(%rsp), %rdx
14008ce57: e8 74 c2 ff ff              	callq	0x1400890d0 <.text+0x880d0>
14008ce5c: 4c 8b 84 24 68 02 00 00     	movq	0x268(%rsp), %r8
14008ce64: 48 8b 8c 24 60 02 00 00     	movq	0x260(%rsp), %rcx
14008ce6c: e9 5a fb ff ff              	jmp	0x14008c9cb <.text+0x8b9cb>
14008ce71: 48 8b 4c 24 48              	movq	0x48(%rsp), %rcx
14008ce76: 45 31 c0                    	xorl	%r8d, %r8d
14008ce79: 41 b9 0a 00 00 00           	movl	$0xa, %r9d
14008ce7f: 48 8d 15 ff df 06 00        	leaq	0x6dfff(%rip), %rdx     # 0x1400fae85
14008ce86: e8 55 a5 02 00              	callq	0x1400b73e0 <.text+0xb63e0>
14008ce8b: 48 83 f8 ff                 	cmpq	$-0x1, %rax
14008ce8f: 0f 85 79 ff ff ff           	jne	0x14008ce0e <.text+0x8be0e>
14008ce95: 4c 8d 0d 72 df 06 00        	leaq	0x6df72(%rip), %r9      # 0x1400fae0e
14008ce9c: 48 83 ff ff                 	cmpq	$-0x1, %rdi
14008cea0: 0f 84 79 ff ff ff           	je	0x14008ce1f <.text+0x8be1f>
14008cea6: c6 03 01                    	movb	$0x1, (%rbx)
14008cea9: 48 8b 8c 24 80 01 00 00     	movq	0x180(%rsp), %rcx
14008ceb1: 41 b8 0a 00 00 00           	movl	$0xa, %r8d
14008ceb7: 31 d2                       	xorl	%edx, %edx
14008ceb9: ff 15 b1 91 0e 00           	callq	*0xe91b1(%rip)          # 0x140176070
14008cebf: 48 89 43 28                 	movq	%rax, 0x28(%rbx)
14008cec3: 48 8b 8c 24 a0 01 00 00     	movq	0x1a0(%rsp), %rcx
14008cecb: 41 b8 0a 00 00 00           	movl	$0xa, %r8d
14008ced1: 31 d2                       	xorl	%edx, %edx
14008ced3: ff 15 8f 91 0e 00           	callq	*0xe918f(%rip)          # 0x140176068
14008ced9: 48 89 43 30                 	movq	%rax, 0x30(%rbx)
14008cedd: 48 8b 4c 24 30              	movq	0x30(%rsp), %rcx
14008cee2: 48 8d 94 24 c0 01 00 00     	leaq	0x1c0(%rsp), %rdx
14008ceea: e8 e1 c1 ff ff              	callq	0x1400890d0 <.text+0x880d0>
14008ceef: 48 8b 4b 40                 	movq	0x40(%rbx), %rcx
14008cef3: 48 8b 84 24 60 02 00 00     	movq	0x260(%rsp), %rax
14008cefb: 4c 8b 84 24 68 02 00 00     	movq	0x268(%rsp), %r8
14008cf03: 49 89 c9                    	movq	%rcx, %r9
14008cf06: 49 39 ce                    	cmpq	%rcx, %r14
14008cf09: 0f 84 4a 05 00 00           	je	0x14008d459 <.text+0x8c459>
14008cf0f: 4c 39 f8                    	cmpq	%r15, %rax
14008cf12: 0f 84 01 05 00 00           	je	0x14008d419 <.text+0x8c419>
14008cf18: 48 8b 53 50                 	movq	0x50(%rbx), %rdx
14008cf1c: 66 49 0f 6e c0              	movq	%r8, %xmm0
14008cf21: 48 89 43 40                 	movq	%rax, 0x40(%rbx)
14008cf25: 0f 16 84 24 70 02 00 00     	movhps	0x270(%rsp), %xmm0      # xmm0 = xmm0[0,1],mem[0,1]
14008cf2d: 0f 11 43 48                 	movups	%xmm0, 0x48(%rbx)
14008cf31: 48 85 c9                    	testq	%rcx, %rcx
14008cf34: 74 26                       	je	0x14008cf5c <.text+0x8bf5c>
14008cf36: 48 89 94 24 70 02 00 00     	movq	%rdx, 0x270(%rsp)
14008cf3e: 45 31 c9                    	xorl	%r9d, %r9d
14008cf41: 66 44 89 09                 	movw	%r9w, (%rcx)
14008cf45: 4c 39 f9                    	cmpq	%r15, %rcx
14008cf48: 74 12                       	je	0x14008cf5c <.text+0x8bf5c>
14008cf4a: 48 8b 84 24 70 02 00 00     	movq	0x270(%rsp), %rax
14008cf52: 48 8d 54 00 02              	leaq	0x2(%rax,%rax), %rdx
14008cf57: e8 74 d1 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008cf5c: 48 8b 4c 24 30              	movq	0x30(%rsp), %rcx
14008cf61: 48 8d 94 24 e0 01 00 00     	leaq	0x1e0(%rsp), %rdx
14008cf69: e8 62 c1 ff ff              	callq	0x1400890d0 <.text+0x880d0>
14008cf6e: 48 8b 4b 60                 	movq	0x60(%rbx), %rcx
14008cf72: 48 8b bc 24 68 02 00 00     	movq	0x268(%rsp), %rdi
14008cf7a: 48 8b 84 24 60 02 00 00     	movq	0x260(%rsp), %rax
14008cf82: 49 89 c9                    	movq	%rcx, %r9
14008cf85: 49 39 cd                    	cmpq	%rcx, %r13
14008cf88: 0f 84 63 03 00 00           	je	0x14008d2f1 <.text+0x8c2f1>
14008cf8e: 4c 39 f8                    	cmpq	%r15, %rax
14008cf91: 0f 84 15 03 00 00           	je	0x14008d2ac <.text+0x8c2ac>
14008cf97: 48 8b 53 70                 	movq	0x70(%rbx), %rdx
14008cf9b: 66 48 0f 6e c7              	movq	%rdi, %xmm0
14008cfa0: 48 89 43 60                 	movq	%rax, 0x60(%rbx)
14008cfa4: 0f 16 84 24 70 02 00 00     	movhps	0x270(%rsp), %xmm0      # xmm0 = xmm0[0,1],mem[0,1]
14008cfac: 0f 11 43 68                 	movups	%xmm0, 0x68(%rbx)
14008cfb0: 48 85 c9                    	testq	%rcx, %rcx
14008cfb3: 74 2a                       	je	0x14008cfdf <.text+0x8bfdf>
14008cfb5: 48 89 94 24 70 02 00 00     	movq	%rdx, 0x270(%rsp)
14008cfbd: 45 31 c0                    	xorl	%r8d, %r8d
14008cfc0: 66 44 89 01                 	movw	%r8w, (%rcx)
14008cfc4: 4c 39 f9                    	cmpq	%r15, %rcx
14008cfc7: 74 12                       	je	0x14008cfdb <.text+0x8bfdb>
14008cfc9: 48 8b 84 24 70 02 00 00     	movq	0x270(%rsp), %rax
14008cfd1: 48 8d 54 00 02              	leaq	0x2(%rax,%rax), %rdx
14008cfd6: e8 f5 d0 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008cfdb: 48 8b 7b 68                 	movq	0x68(%rbx), %rdi
14008cfdf: 4c 8d 2d 4e dd 06 00        	leaq	0x6dd4e(%rip), %r13     # 0x1400fad34
14008cfe6: 4c 89 e9                    	movq	%r13, %rcx
14008cfe9: e8 92 86 01 00              	callq	0x1400a5680 <.text+0xa4680>
14008cfee: 48 39 c7                    	cmpq	%rax, %rdi
14008cff1: 0f 84 2a 02 00 00           	je	0x14008d221 <.text+0x8c221>
14008cff7: 4c 8b 73 68                 	movq	0x68(%rbx), %r14
14008cffb: 48 8d 3d 4a dd 06 00        	leaq	0x6dd4a(%rip), %rdi     # 0x1400fad4c
14008d002: 48 89 f9                    	movq	%rdi, %rcx
14008d005: e8 76 86 01 00              	callq	0x1400a5680 <.text+0xa4680>
14008d00a: 49 39 c6                    	cmpq	%rax, %r14
14008d00d: 0f 84 5d 02 00 00           	je	0x14008d270 <.text+0x8c270>
14008d013: 48 8b 4c 24 30              	movq	0x30(%rsp), %rcx
14008d018: 48 8d 94 24 00 02 00 00     	leaq	0x200(%rsp), %rdx
14008d020: e8 ab c0 ff ff              	callq	0x1400890d0 <.text+0x880d0>
14008d025: 48 8b 8b 80 00 00 00        	movq	0x80(%rbx), %rcx
14008d02c: 48 8b 84 24 60 02 00 00     	movq	0x260(%rsp), %rax
14008d034: 4c 8b 84 24 68 02 00 00     	movq	0x268(%rsp), %r8
14008d03c: 49 89 c9                    	movq	%rcx, %r9
14008d03f: 49 39 cc                    	cmpq	%rcx, %r12
14008d042: 0f 84 ac 03 00 00           	je	0x14008d3f4 <.text+0x8c3f4>
14008d048: 4c 39 f8                    	cmpq	%r15, %rax
14008d04b: 0f 84 59 03 00 00           	je	0x14008d3aa <.text+0x8c3aa>
14008d051: 48 8b 93 90 00 00 00        	movq	0x90(%rbx), %rdx
14008d058: 66 49 0f 6e c0              	movq	%r8, %xmm0
14008d05d: 48 89 83 80 00 00 00        	movq	%rax, 0x80(%rbx)
14008d064: 0f 16 84 24 70 02 00 00     	movhps	0x270(%rsp), %xmm0      # xmm0 = xmm0[0,1],mem[0,1]
14008d06c: 0f 11 83 88 00 00 00        	movups	%xmm0, 0x88(%rbx)
14008d073: 48 85 c9                    	testq	%rcx, %rcx
14008d076: 74 24                       	je	0x14008d09c <.text+0x8c09c>
14008d078: 48 89 94 24 70 02 00 00     	movq	%rdx, 0x270(%rsp)
14008d080: 31 c0                       	xorl	%eax, %eax
14008d082: 66 89 01                    	movw	%ax, (%rcx)
14008d085: 4c 39 f9                    	cmpq	%r15, %rcx
14008d088: 74 12                       	je	0x14008d09c <.text+0x8c09c>
14008d08a: 48 8b 84 24 70 02 00 00     	movq	0x270(%rsp), %rax
14008d092: 48 8d 54 00 02              	leaq	0x2(%rax,%rax), %rdx
14008d097: e8 34 d0 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008d09c: 48 8b 83 a0 00 00 00        	movq	0xa0(%rbx), %rax
14008d0a3: 48 89 c2                    	movq	%rax, %rdx
14008d0a6: 48 39 c6                    	cmpq	%rax, %rsi
14008d0a9: 0f 84 61 02 00 00           	je	0x14008d310 <.text+0x8c310>
14008d0af: 48 8b 8c 24 20 02 00 00     	movq	0x220(%rsp), %rcx
14008d0b7: 48 3b 4c 24 38              	cmpq	0x38(%rsp), %rcx
14008d0bc: 0f 84 9e 02 00 00           	je	0x14008d360 <.text+0x8c360>
14008d0c2: 4c 8b 83 b0 00 00 00        	movq	0xb0(%rbx), %r8
14008d0c9: 48 89 8b a0 00 00 00        	movq	%rcx, 0xa0(%rbx)
14008d0d0: f3 0f 7e 84 24 28 02 00 00  	movq	0x228(%rsp), %xmm0
14008d0d9: 0f 16 84 24 30 02 00 00     	movhps	0x230(%rsp), %xmm0      # xmm0 = xmm0[0,1],mem[0,1]
14008d0e1: 0f 11 83 a8 00 00 00        	movups	%xmm0, 0xa8(%rbx)
14008d0e8: 48 85 c0                    	testq	%rax, %rax
14008d0eb: 0f 84 4d 02 00 00           	je	0x14008d33e <.text+0x8c33e>
14008d0f1: 48 89 84 24 20 02 00 00     	movq	%rax, 0x220(%rsp)
14008d0f9: 4c 89 84 24 30 02 00 00     	movq	%r8, 0x230(%rsp)
14008d101: c6 02 00                    	movb	$0x0, (%rdx)
14008d104: 48 8d 0d 85 dd 06 00        	leaq	0x6dd85(%rip), %rcx     # 0x1400fae90
14008d10b: 48 c7 84 24 28 02 00 00 00 00 00 00 	movq	$0x0, 0x228(%rsp)
14008d117: e8 64 85 01 00              	callq	0x1400a5680 <.text+0xa4680>
14008d11c: 48 8d 4b 08                 	leaq	0x8(%rbx), %rcx
14008d120: 4c 8d 0d 69 dd 06 00        	leaq	0x6dd69(%rip), %r9      # 0x1400fae90
14008d127: 31 d2                       	xorl	%edx, %edx
14008d129: 48 89 44 24 20              	movq	%rax, 0x20(%rsp)
14008d12e: 4c 8b 43 10                 	movq	0x10(%rbx), %r8
14008d132: e8 e9 6b 03 00              	callq	0x1400c3d20 <.text+0xc2d20>
14008d137: e9 cb fa ff ff              	jmp	0x14008cc07 <.text+0x8bc07>
14008d13c: 4d 85 c0                    	testq	%r8, %r8
14008d13f: 0f 84 d1 00 00 00           	je	0x14008d216 <.text+0x8c216>
14008d145: 49 83 f8 01                 	cmpq	$0x1, %r8
14008d149: 0f 84 bc 00 00 00           	je	0x14008d20b <.text+0x8c20b>
14008d14f: 4d 01 c0                    	addq	%r8, %r8
14008d152: 4c 89 fa                    	movq	%r15, %rdx
14008d155: 48 89 c1                    	movq	%rax, %rcx
14008d158: e8 8b 86 01 00              	callq	0x1400a57e8 <.text+0xa47e8>
14008d15d: 4c 8b 84 24 68 02 00 00     	movq	0x268(%rsp), %r8
14008d165: 48 8b 53 08                 	movq	0x8(%rbx), %rdx
14008d169: 48 8b 84 24 60 02 00 00     	movq	0x260(%rsp), %rax
14008d171: 31 f6                       	xorl	%esi, %esi
14008d173: 66 42 89 34 42              	movw	%si, (%rdx,%r8,2)
14008d178: 4c 89 43 10                 	movq	%r8, 0x10(%rbx)
14008d17c: e9 85 f8 ff ff              	jmp	0x14008ca06 <.text+0x8ba06>
14008d181: 4c 39 f9                    	cmpq	%r15, %rcx
14008d184: 74 b6                       	je	0x14008d13c <.text+0x8c13c>
14008d186: 48 8b 84 24 70 02 00 00     	movq	0x270(%rsp), %rax
14008d18e: 48 89 4b 08                 	movq	%rcx, 0x8(%rbx)
14008d192: 4c 89 43 10                 	movq	%r8, 0x10(%rbx)
14008d196: 48 89 43 18                 	movq	%rax, 0x18(%rbx)
14008d19a: e9 88 f8 ff ff              	jmp	0x14008ca27 <.text+0x8ba27>
14008d19f: 41 bd 51 51 00 00           	movl	$0x5151, %r13d          # imm = 0x5151
14008d1a5: b8 02 00 00 00              	movl	$0x2, %eax
14008d1aa: 66 44 89 29                 	movw	%r13w, (%rcx)
14008d1ae: e9 08 f8 ff ff              	jmp	0x14008c9bb <.text+0x8b9bb>
14008d1b3: 48 8b 4c 24 30              	movq	0x30(%rsp), %rcx
14008d1b8: 48 8d 94 24 20 02 00 00     	leaq	0x220(%rsp), %rdx
14008d1c0: 45 31 c0                    	xorl	%r8d, %r8d
14008d1c3: e8 18 96 03 00              	callq	0x1400c67e0 <.text+0xc57e0>
14008d1c8: 48 89 84 24 60 02 00 00     	movq	%rax, 0x260(%rsp)
14008d1d0: 48 89 c1                    	movq	%rax, %rcx
14008d1d3: 48 8b 84 24 20 02 00 00     	movq	0x220(%rsp), %rax
14008d1db: 48 89 84 24 70 02 00 00     	movq	%rax, 0x270(%rsp)
14008d1e3: 49 89 f0                    	movq	%rsi, %r8
14008d1e6: 48 8d 15 61 dc 06 00        	leaq	0x6dc61(%rip), %rdx     # 0x1400fae4e
14008d1ed: e8 f6 85 01 00              	callq	0x1400a57e8 <.text+0xa47e8>
14008d1f2: 4c 8b 84 24 20 02 00 00     	movq	0x220(%rsp), %r8
14008d1fa: 48 8b 8c 24 60 02 00 00     	movq	0x260(%rsp), %rcx
14008d202: 4b 8d 04 00                 	leaq	(%r8,%r8), %rax
14008d206: e9 b0 f7 ff ff              	jmp	0x14008c9bb <.text+0x8b9bb>
14008d20b: 0f b7 94 24 70 02 00 00     	movzwl	0x270(%rsp), %edx
14008d213: 66 89 10                    	movw	%dx, (%rax)
14008d216: 48 89 c2                    	movq	%rax, %rdx
14008d219: 48 89 c8                    	movq	%rcx, %rax
14008d21c: e9 50 ff ff ff              	jmp	0x14008d171 <.text+0x8c171>
14008d221: 4c 8b 73 68                 	movq	0x68(%rbx), %r14
14008d225: 48 8b 43 60                 	movq	0x60(%rbx), %rax
14008d229: 4d 85 f6                    	testq	%r14, %r14
14008d22c: 74 30                       	je	0x14008d25e <.text+0x8c25e>
14008d22e: 4c 39 e8                    	cmpq	%r13, %rax
14008d231: 74 2b                       	je	0x14008d25e <.text+0x8c25e>
14008d233: 48 85 c0                    	testq	%rax, %rax
14008d236: 0f 84 bf fd ff ff           	je	0x14008cffb <.text+0x8bffb>
14008d23c: 4c 89 f1                    	movq	%r14, %rcx
14008d23f: 48 89 c2                    	movq	%rax, %rdx
14008d242: 41 0f b7 7d 00              	movzwl	(%r13), %edi
14008d247: 66 39 3a                    	cmpw	%di, (%rdx)
14008d24a: 0f 85 ab fd ff ff           	jne	0x14008cffb <.text+0x8bffb>
14008d250: 48 83 c2 02                 	addq	$0x2, %rdx
14008d254: 49 83 c5 02                 	addq	$0x2, %r13
14008d258: 48 83 e9 01                 	subq	$0x1, %rcx
14008d25c: 75 e4                       	jne	0x14008d242 <.text+0x8c242>
14008d25e: 31 d2                       	xorl	%edx, %edx
14008d260: 66 89 10                    	movw	%dx, (%rax)
14008d263: 48 c7 43 68 00 00 00 00     	movq	$0x0, 0x68(%rbx)
14008d26b: e9 a3 fd ff ff              	jmp	0x14008d013 <.text+0x8c013>
14008d270: 48 8b 43 60                 	movq	0x60(%rbx), %rax
14008d274: 48 8b 53 68                 	movq	0x68(%rbx), %rdx
14008d278: 48 39 f8                    	cmpq	%rdi, %rax
14008d27b: 74 e1                       	je	0x14008d25e <.text+0x8c25e>
14008d27d: 48 85 d2                    	testq	%rdx, %rdx
14008d280: 74 dc                       	je	0x14008d25e <.text+0x8c25e>
14008d282: 48 85 c0                    	testq	%rax, %rax
14008d285: 0f 84 88 fd ff ff           	je	0x14008d013 <.text+0x8c013>
14008d28b: 48 89 c1                    	movq	%rax, %rcx
14008d28e: 44 0f b7 1f                 	movzwl	(%rdi), %r11d
14008d292: 66 44 39 19                 	cmpw	%r11w, (%rcx)
14008d296: 0f 85 77 fd ff ff           	jne	0x14008d013 <.text+0x8c013>
14008d29c: 48 83 c1 02                 	addq	$0x2, %rcx
14008d2a0: 48 83 c7 02                 	addq	$0x2, %rdi
14008d2a4: 48 83 ea 01                 	subq	$0x1, %rdx
14008d2a8: 75 e4                       	jne	0x14008d28e <.text+0x8c28e>
14008d2aa: eb b2                       	jmp	0x14008d25e <.text+0x8c25e>
14008d2ac: 48 85 ff                    	testq	%rdi, %rdi
14008d2af: 74 2d                       	je	0x14008d2de <.text+0x8c2de>
14008d2b1: 48 83 ff 01                 	cmpq	$0x1, %rdi
14008d2b5: 0f 84 cb 01 00 00           	je	0x14008d486 <.text+0x8c486>
14008d2bb: 4c 8d 04 3f                 	leaq	(%rdi,%rdi), %r8
14008d2bf: 4c 89 c9                    	movq	%r9, %rcx
14008d2c2: 4c 89 fa                    	movq	%r15, %rdx
14008d2c5: e8 1e 85 01 00              	callq	0x1400a57e8 <.text+0xa47e8>
14008d2ca: 48 8b bc 24 68 02 00 00     	movq	0x268(%rsp), %rdi
14008d2d2: 4c 8b 4b 60                 	movq	0x60(%rbx), %r9
14008d2d6: 48 8b 84 24 60 02 00 00     	movq	0x260(%rsp), %rax
14008d2de: 66 41 c7 04 79 00 00        	movw	$0x0, (%r9,%rdi,2)
14008d2e5: 48 89 c1                    	movq	%rax, %rcx
14008d2e8: 48 89 7b 68                 	movq	%rdi, 0x68(%rbx)
14008d2ec: e9 cc fc ff ff              	jmp	0x14008cfbd <.text+0x8bfbd>
14008d2f1: 4c 39 f8                    	cmpq	%r15, %rax
14008d2f4: 74 b6                       	je	0x14008d2ac <.text+0x8c2ac>
14008d2f6: 66 48 0f 6e c7              	movq	%rdi, %xmm0
14008d2fb: 48 89 43 60                 	movq	%rax, 0x60(%rbx)
14008d2ff: 0f 16 84 24 70 02 00 00     	movhps	0x270(%rsp), %xmm0      # xmm0 = xmm0[0,1],mem[0,1]
14008d307: 0f 11 43 68                 	movups	%xmm0, 0x68(%rbx)
14008d30b: e9 cf fc ff ff              	jmp	0x14008cfdf <.text+0x8bfdf>
14008d310: 48 8b 94 24 20 02 00 00     	movq	0x220(%rsp), %rdx
14008d318: 48 3b 54 24 38              	cmpq	0x38(%rsp), %rdx
14008d31d: 74 41                       	je	0x14008d360 <.text+0x8c360>
14008d31f: 48 89 93 a0 00 00 00        	movq	%rdx, 0xa0(%rbx)
14008d326: f3 0f 7e 84 24 28 02 00 00  	movq	0x228(%rsp), %xmm0
14008d32f: 0f 16 84 24 30 02 00 00     	movhps	0x230(%rsp), %xmm0      # xmm0 = xmm0[0,1],mem[0,1]
14008d337: 0f 11 83 a8 00 00 00        	movups	%xmm0, 0xa8(%rbx)
14008d33e: 48 8b 44 24 38              	movq	0x38(%rsp), %rax
14008d343: 48 89 84 24 20 02 00 00     	movq	%rax, 0x220(%rsp)
14008d34b: 48 8d 84 24 30 02 00 00     	leaq	0x230(%rsp), %rax
14008d353: 48 89 44 24 38              	movq	%rax, 0x38(%rsp)
14008d358: 48 89 c2                    	movq	%rax, %rdx
14008d35b: e9 a1 fd ff ff              	jmp	0x14008d101 <.text+0x8c101>
14008d360: 48 8b 94 24 28 02 00 00     	movq	0x228(%rsp), %rdx
14008d368: 48 85 d2                    	testq	%rdx, %rdx
14008d36b: 74 25                       	je	0x14008d392 <.text+0x8c392>
14008d36d: 48 83 fa 01                 	cmpq	$0x1, %rdx
14008d371: 0f 84 31 01 00 00           	je	0x14008d4a8 <.text+0x8c4a8>
14008d377: 48 8b 74 24 38              	movq	0x38(%rsp), %rsi
14008d37c: 89 d1                       	movl	%edx, %ecx
14008d37e: 48 89 c7                    	movq	%rax, %rdi
14008d381: f3 a4                       	rep		movsb	(%rsi), %es:(%rdi)
14008d383: 48 8b 94 24 28 02 00 00     	movq	0x228(%rsp), %rdx
14008d38b: 48 8b 83 a0 00 00 00        	movq	0xa0(%rbx), %rax
14008d392: c6 04 10 00                 	movb	$0x0, (%rax,%rdx)
14008d396: 48 89 93 a8 00 00 00        	movq	%rdx, 0xa8(%rbx)
14008d39d: 48 8b 94 24 20 02 00 00     	movq	0x220(%rsp), %rdx
14008d3a5: e9 57 fd ff ff              	jmp	0x14008d101 <.text+0x8c101>
14008d3aa: 4d 85 c0                    	testq	%r8, %r8
14008d3ad: 74 2f                       	je	0x14008d3de <.text+0x8c3de>
14008d3af: 49 83 f8 01                 	cmpq	$0x1, %r8
14008d3b3: 0f 84 de 00 00 00           	je	0x14008d497 <.text+0x8c497>
14008d3b9: 4d 01 c0                    	addq	%r8, %r8
14008d3bc: 4c 89 c9                    	movq	%r9, %rcx
14008d3bf: 4c 89 fa                    	movq	%r15, %rdx
14008d3c2: e8 21 84 01 00              	callq	0x1400a57e8 <.text+0xa47e8>
14008d3c7: 4c 8b 84 24 68 02 00 00     	movq	0x268(%rsp), %r8
14008d3cf: 48 8b 84 24 60 02 00 00     	movq	0x260(%rsp), %rax
14008d3d7: 4c 8b 8b 80 00 00 00        	movq	0x80(%rbx), %r9
14008d3de: 66 43 c7 04 41 00 00        	movw	$0x0, (%r9,%r8,2)
14008d3e5: 48 89 c1                    	movq	%rax, %rcx
14008d3e8: 4c 89 83 88 00 00 00        	movq	%r8, 0x88(%rbx)
14008d3ef: e9 8c fc ff ff              	jmp	0x14008d080 <.text+0x8c080>
14008d3f4: 4c 39 f8                    	cmpq	%r15, %rax
14008d3f7: 74 b1                       	je	0x14008d3aa <.text+0x8c3aa>
14008d3f9: 66 49 0f 6e c0              	movq	%r8, %xmm0
14008d3fe: 48 89 83 80 00 00 00        	movq	%rax, 0x80(%rbx)
