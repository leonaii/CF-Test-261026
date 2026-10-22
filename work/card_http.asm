
CF上号器_v1.17.exe:	file format coff-x86-64

Disassembly of section .text:

0000000140001000 <.text>:
1400896f0: 41 57                       	pushq	%r15
1400896f2: 41 56                       	pushq	%r14
1400896f4: 41 55                       	pushq	%r13
1400896f6: 41 54                       	pushq	%r12
1400896f8: 55                          	pushq	%rbp
1400896f9: 57                          	pushq	%rdi
1400896fa: 56                          	pushq	%rsi
1400896fb: 53                          	pushq	%rbx
1400896fc: 48 81 ec 38 0b 00 00        	subq	$0xb38, %rsp            # imm = 0xB38
140089703: 66 0f ef c0                 	pxor	%xmm0, %xmm0
140089707: 49 8b 01                    	movq	(%r9), %rax
14008970a: c6 00 00                    	movb	$0x0, (%rax)
14008970d: 48 89 cd                    	movq	%rcx, %rbp
140089710: 48 89 d6                    	movq	%rdx, %rsi
140089713: 49 c7 41 08 00 00 00 00     	movq	$0x0, 0x8(%r9)
14008971b: 31 d2                       	xorl	%edx, %edx
14008971d: 4c 8d ac 24 30 01 00 00     	leaq	0x130(%rsp), %r13
140089725: 4c 8d b4 24 30 03 00 00     	leaq	0x330(%rsp), %r14
14008972d: 4c 89 c7                    	movq	%r8, %rdi
140089730: 41 b8 00 02 00 00           	movl	$0x200, %r8d            # imm = 0x200
140089736: 4c 89 e9                    	movq	%r13, %rcx
140089739: 4c 89 cb                    	movq	%r9, %rbx
14008973c: 0f 11 84 24 c4 00 00 00     	movups	%xmm0, 0xc4(%rsp)
140089744: 0f 11 84 24 d4 00 00 00     	movups	%xmm0, 0xd4(%rsp)
14008974c: 0f 11 84 24 e4 00 00 00     	movups	%xmm0, 0xe4(%rsp)
140089754: 0f 11 84 24 f4 00 00 00     	movups	%xmm0, 0xf4(%rsp)
14008975c: 0f 11 84 24 04 01 00 00     	movups	%xmm0, 0x104(%rsp)
140089764: 0f 11 84 24 14 01 00 00     	movups	%xmm0, 0x114(%rsp)
14008976c: c7 84 24 24 01 00 00 00 00 00 00    	movl	$0x0, 0x124(%rsp)
140089777: c7 84 24 c0 00 00 00 68 00 00 00    	movl	$0x68, 0xc0(%rsp)
140089782: e8 c9 be 01 00              	callq	0x1400a5650 <.text+0xa4650>
140089787: 41 b8 00 08 00 00           	movl	$0x800, %r8d            # imm = 0x800
14008978d: 31 d2                       	xorl	%edx, %edx
14008978f: 4c 89 f1                    	movq	%r14, %rcx
140089792: e8 b9 be 01 00              	callq	0x1400a5650 <.text+0xa4650>
140089797: 4c 8b 45 08                 	movq	0x8(%rbp), %r8
14008979b: 4c 8b 26                    	movq	(%rsi), %r12
14008979e: 45 31 d2                    	xorl	%r10d, %r10d
1400897a1: 48 8b 76 08                 	movq	0x8(%rsi), %rsi
1400897a5: 4c 8b 7d 00                 	movq	(%rbp), %r15
1400897a9: 48 8d 4c 24 60              	leaq	0x60(%rsp), %rcx
1400897ae: 48 8d 6c 24 70              	leaq	0x70(%rsp), %rbp
1400897b3: 4c 89 ac 24 d8 00 00 00     	movq	%r13, 0xd8(%rsp)
1400897bb: 49 8d 14 30                 	leaq	(%r8,%rsi), %rdx
1400897bf: 48 89 6c 24 60              	movq	%rbp, 0x60(%rsp)
1400897c4: c7 84 24 e0 00 00 00 00 01 00 00    	movl	$0x100, 0xe0(%rsp) # imm = 0x100
1400897cf: 4c 89 b4 24 08 01 00 00     	movq	%r14, 0x108(%rsp)
1400897d7: c7 84 24 10 01 00 00 00 04 00 00    	movl	$0x400, 0x110(%rsp) # imm = 0x400
1400897e2: 66 44 89 54 24 70           	movw	%r10w, 0x70(%rsp)
1400897e8: 48 c7 44 24 68 00 00 00 00  	movq	$0x0, 0x68(%rsp)
1400897f1: 4c 89 44 24 48              	movq	%r8, 0x48(%rsp)
1400897f6: e8 c5 cc 03 00              	callq	0x1400c64c0 <.text+0xc54c0>
1400897fb: 48 b8 fe ff ff ff ff ff ff 3f       	movabsq	$0x3ffffffffffffffe, %rax # imm = 0x3FFFFFFFFFFFFFFE
140089805: 4c 8b 44 24 48              	movq	0x48(%rsp), %r8
14008980a: 48 2b 44 24 68              	subq	0x68(%rsp), %rax
14008980f: 4c 39 c0                    	cmpq	%r8, %rax
140089812: 0f 82 e3 e8 04 00           	jb	0x1400d80fb <.text+0xd70fb>
140089818: 4c 89 fa                    	movq	%r15, %rdx
14008981b: 48 8d 4c 24 60              	leaq	0x60(%rsp), %rcx
140089820: e8 db cd 03 00              	callq	0x1400c6600 <.text+0xc5600>
140089825: 48 b8 fe ff ff ff ff ff ff 3f       	movabsq	$0x3ffffffffffffffe, %rax # imm = 0x3FFFFFFFFFFFFFFE
14008982f: 48 2b 44 24 68              	subq	0x68(%rsp), %rax
140089834: 48 39 f0                    	cmpq	%rsi, %rax
140089837: 0f 82 b2 e8 04 00           	jb	0x1400d80ef <.text+0xd70ef>
14008983d: 49 89 f0                    	movq	%rsi, %r8
140089840: 4c 89 e2                    	movq	%r12, %rdx
140089843: 48 8d 4c 24 60              	leaq	0x60(%rsp), %rcx
140089848: e8 b3 cd 03 00              	callq	0x1400c6600 <.text+0xc5600>
14008984d: 48 8b 4c 24 60              	movq	0x60(%rsp), %rcx
140089852: 4c 8d 8c 24 c0 00 00 00     	leaq	0xc0(%rsp), %r9
14008985a: 45 31 c0                    	xorl	%r8d, %r8d
14008985d: 31 d2                       	xorl	%edx, %edx
14008985f: ff 15 63 cc 0e 00           	callq	*0xecc63(%rip)          # 0x1401764c8
140089865: 85 c0                       	testl	%eax, %eax
140089867: 0f 84 43 02 00 00           	je	0x140089ab0 <.text+0x88ab0>
14008986d: c7 44 24 20 00 00 00 00     	movl	$0x0, 0x20(%rsp)
140089875: 45 31 c9                    	xorl	%r9d, %r9d
140089878: 45 31 c0                    	xorl	%r8d, %r8d
14008987b: 31 d2                       	xorl	%edx, %edx
14008987d: 48 8d 0d 9e 13 07 00        	leaq	0x7139e(%rip), %rcx     # 0x1400fac22
140089884: ff 15 46 cc 0e 00           	callq	*0xecc46(%rip)          # 0x1401764d0
14008988a: 49 89 c4                    	movq	%rax, %r12
14008988d: 48 85 c0                    	testq	%rax, %rax
140089890: 0f 84 e2 03 00 00           	je	0x140089c78 <.text+0x88c78>
140089896: 44 0f b7 84 24 e4 00 00 00  	movzwl	0xe4(%rsp), %r8d
14008989f: 45 31 c9                    	xorl	%r9d, %r9d
1400898a2: 4c 89 ea                    	movq	%r13, %rdx
1400898a5: 48 89 c1                    	movq	%rax, %rcx
1400898a8: ff 15 12 cc 0e 00           	callq	*0xecc12(%rip)          # 0x1401764c0
1400898ae: 49 89 c5                    	movq	%rax, %r13
1400898b1: 48 85 c0                    	testq	%rax, %rax
1400898b4: 0f 84 66 02 00 00           	je	0x140089b20 <.text+0x88b20>
1400898ba: 31 c0                       	xorl	%eax, %eax
1400898bc: 83 bc 24 d4 00 00 00 02     	cmpl	$0x2, 0xd4(%rsp)
1400898c4: 4d 89 f0                    	movq	%r14, %r8
1400898c7: 4c 89 e9                    	movq	%r13, %rcx
1400898ca: 0f 94 c0                    	sete	%al
1400898cd: 48 8d 15 92 13 07 00        	leaq	0x71392(%rip), %rdx     # 0x1400fac66
1400898d4: 45 31 c9                    	xorl	%r9d, %r9d
1400898d7: 48 c7 44 24 28 00 00 00 00  	movq	$0x0, 0x28(%rsp)
1400898e0: 48 c7 44 24 20 00 00 00 00  	movq	$0x0, 0x20(%rsp)
1400898e9: c1 e0 17                    	shll	$0x17, %eax
1400898ec: 89 44 24 30                 	movl	%eax, 0x30(%rsp)
1400898f0: ff 15 e2 cb 0e 00           	callq	*0xecbe2(%rip)          # 0x1401764d8
1400898f6: 48 89 c6                    	movq	%rax, %rsi
1400898f9: 48 85 c0                    	testq	%rax, %rax
1400898fc: 0f 84 96 03 00 00           	je	0x140089c98 <.text+0x88c98>
140089902: c7 44 24 20 98 3a 00 00     	movl	$0x3a98, 0x20(%rsp)     # imm = 0x3A98
14008990a: 41 b9 98 3a 00 00           	movl	$0x3a98, %r9d           # imm = 0x3A98
140089910: ba 88 13 00 00              	movl	$0x1388, %edx           # imm = 0x1388
140089915: 48 89 c1                    	movq	%rax, %rcx
140089918: 41 b8 88 13 00 00           	movl	$0x1388, %r8d           # imm = 0x1388
14008991e: ff 15 ec cb 0e 00           	callq	*0xecbec(%rip)          # 0x140176510
140089924: 48 8d 84 24 90 00 00 00     	leaq	0x90(%rsp), %rax
14008992c: 48 8d 0d 4d 13 07 00        	leaq	0x7134d(%rip), %rcx     # 0x1400fac80
140089933: 48 89 84 24 80 00 00 00     	movq	%rax, 0x80(%rsp)
14008993b: e8 40 bd 01 00              	callq	0x1400a5680 <.text+0xa4680>
140089940: 4c 8d 34 00                 	leaq	(%rax,%rax), %r14
140089944: 4c 89 f0                    	movq	%r14, %rax
140089947: 48 d1 f8                    	sarq	%rax
14008994a: 48 89 84 24 a0 00 00 00     	movq	%rax, 0xa0(%rsp)
140089952: 49 83 fe 0e                 	cmpq	$0xe, %r14
140089956: 0f 87 bc 02 00 00           	ja	0x140089c18 <.text+0x88c18>
14008995c: 48 8b 8c 24 80 00 00 00     	movq	0x80(%rsp), %rcx
140089964: 48 83 f8 01                 	cmpq	$0x1, %rax
140089968: 0f 85 a2 01 00 00           	jne	0x140089b10 <.text+0x88b10>
14008996e: 41 b9 43 00 00 00           	movl	$0x43, %r9d
140089974: ba 02 00 00 00              	movl	$0x2, %edx
140089979: 66 44 89 09                 	movw	%r9w, (%rcx)
14008997d: 45 31 c0                    	xorl	%r8d, %r8d
140089980: 66 44 89 04 11              	movw	%r8w, (%rcx,%rdx)
140089985: 48 89 ca                    	movq	%rcx, %rdx
140089988: 41 b8 ff ff ff ff           	movl	$0xffffffff, %r8d       # imm = 0xFFFFFFFF
14008998e: 48 89 f1                    	movq	%rsi, %rcx
140089991: 48 89 84 24 88 00 00 00     	movq	%rax, 0x88(%rsp)
140089999: 48 8b 47 08                 	movq	0x8(%rdi), %rax
14008999d: 48 c7 44 24 30 00 00 00 00  	movq	$0x0, 0x30(%rsp)
1400899a6: 89 44 24 28                 	movl	%eax, 0x28(%rsp)
1400899aa: 89 44 24 20                 	movl	%eax, 0x20(%rsp)
1400899ae: 4c 8b 0f                    	movq	(%rdi), %r9
1400899b1: ff 15 49 cb 0e 00           	callq	*0xecb49(%rip)          # 0x140176500
1400899b7: 85 c0                       	testl	%eax, %eax
1400899b9: 0f 84 89 01 00 00           	je	0x140089b48 <.text+0x88b48>
1400899bf: 31 d2                       	xorl	%edx, %edx
1400899c1: 48 89 f1                    	movq	%rsi, %rcx
1400899c4: ff 15 2e cb 0e 00           	callq	*0xecb2e(%rip)          # 0x1401764f8
1400899ca: 85 c0                       	testl	%eax, %eax
1400899cc: 0f 84 76 01 00 00           	je	0x140089b48 <.text+0x88b48>
1400899d2: 4c 8b 35 07 cb 0e 00        	movq	0xecb07(%rip), %r14     # 0x1401764e0
1400899d9: 4c 8b 3d 10 cb 0e 00        	movq	0xecb10(%rip), %r15     # 0x1401764f0
1400899e0: c7 44 24 58 00 00 00 00     	movl	$0x0, 0x58(%rsp)
1400899e8: 48 8d 54 24 58              	leaq	0x58(%rsp), %rdx
1400899ed: 48 89 f1                    	movq	%rsi, %rcx
1400899f0: 41 ff d6                    	callq	*%r14
1400899f3: 85 c0                       	testl	%eax, %eax
1400899f5: 0f 84 fb 01 00 00           	je	0x140089bf6 <.text+0x88bf6>
1400899fb: 8b 54 24 58                 	movl	0x58(%rsp), %edx
1400899ff: 85 d2                       	testl	%edx, %edx
140089a01: 0f 84 ef 01 00 00           	je	0x140089bf6 <.text+0x88bf6>
140089a07: 48 8d bc 24 b0 00 00 00     	leaq	0xb0(%rsp), %rdi
140089a0f: 45 31 c0                    	xorl	%r8d, %r8d
140089a12: 48 8d 8c 24 a0 00 00 00     	leaq	0xa0(%rsp), %rcx
140089a1a: 48 89 bc 24 a0 00 00 00     	movq	%rdi, 0xa0(%rsp)
140089a22: e8 09 61 03 00              	callq	0x1400bfb30 <.text+0xbeb30>
140089a27: c7 44 24 5c 00 00 00 00     	movl	$0x0, 0x5c(%rsp)
140089a2f: 44 8b 44 24 58              	movl	0x58(%rsp), %r8d
140089a34: 4c 8d 4c 24 5c              	leaq	0x5c(%rsp), %r9
140089a39: 48 89 f1                    	movq	%rsi, %rcx
140089a3c: 48 8b 94 24 a0 00 00 00     	movq	0xa0(%rsp), %rdx
140089a44: 41 ff d7                    	callq	*%r15
140089a47: 48 8b 8c 24 a0 00 00 00     	movq	0xa0(%rsp), %rcx
140089a4f: 85 c0                       	testl	%eax, %eax
140089a51: 0f 84 89 01 00 00           	je	0x140089be0 <.text+0x88be0>
140089a57: 44 8b 44 24 5c              	movl	0x5c(%rsp), %r8d
140089a5c: 45 85 c0                    	testl	%r8d, %r8d
140089a5f: 0f 84 7b 01 00 00           	je	0x140089be0 <.text+0x88be0>
140089a65: 48 b8 fe ff ff ff ff ff ff 7f       	movabsq	$0x7ffffffffffffffe, %rax # imm = 0x7FFFFFFFFFFFFFFE
140089a6f: 48 2b 43 08                 	subq	0x8(%rbx), %rax
140089a73: 4c 39 c0                    	cmpq	%r8, %rax
140089a76: 0f 82 cf e6 04 00           	jb	0x1400d814b <.text+0xd714b>
140089a7c: 48 89 ca                    	movq	%rcx, %rdx
140089a7f: 48 89 d9                    	movq	%rbx, %rcx
140089a82: e8 49 8e 03 00              	callq	0x1400c28d0 <.text+0xc18d0>
140089a87: 48 8b 8c 24 a0 00 00 00     	movq	0xa0(%rsp), %rcx
140089a8f: 48 39 f9                    	cmpq	%rdi, %rcx
140089a92: 0f 84 48 ff ff ff           	je	0x1400899e0 <.text+0x889e0>
140089a98: 48 8b 84 24 b0 00 00 00     	movq	0xb0(%rsp), %rax
140089aa0: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
140089aa4: e8 27 06 04 00              	callq	0x1400ca0d0 <.text+0xc90d0>
140089aa9: e9 32 ff ff ff              	jmp	0x1400899e0 <.text+0x889e0>
140089aae: 66 90                       	nop
140089ab0: 48 8d 0d 5b 11 07 00        	leaq	0x7115b(%rip), %rcx     # 0x1400fac12
140089ab7: e8 c4 bb 01 00              	callq	0x1400a5680 <.text+0xa4680>
140089abc: 4c 8d 0d 4f 11 07 00        	leaq	0x7114f(%rip), %r9      # 0x1400fac12
140089ac3: 48 89 44 24 20              	movq	%rax, 0x20(%rsp)
140089ac8: 48 8b 84 24 a0 0b 00 00     	movq	0xba0(%rsp), %rax
140089ad0: 31 d2                       	xorl	%edx, %edx
140089ad2: 4c 8b 40 08                 	movq	0x8(%rax), %r8
140089ad6: 48 89 c1                    	movq	%rax, %rcx
140089ad9: e8 42 a2 03 00              	callq	0x1400c3d20 <.text+0xc2d20>
140089ade: 31 db                       	xorl	%ebx, %ebx
140089ae0: 48 8b 4c 24 60              	movq	0x60(%rsp), %rcx
140089ae5: 48 39 e9                    	cmpq	%rbp, %rcx
140089ae8: 74 0f                       	je	0x140089af9 <.text+0x88af9>
140089aea: 48 8b 44 24 70              	movq	0x70(%rsp), %rax
140089aef: 48 8d 54 00 02              	leaq	0x2(%rax,%rax), %rdx
140089af4: e8 d7 05 04 00              	callq	0x1400ca0d0 <.text+0xc90d0>
140089af9: 89 d8                       	movl	%ebx, %eax
140089afb: 48 81 c4 38 0b 00 00        	addq	$0xb38, %rsp            # imm = 0xB38
140089b02: 5b                          	popq	%rbx
140089b03: 5e                          	popq	%rsi
140089b04: 5f                          	popq	%rdi
140089b05: 5d                          	popq	%rbp
140089b06: 41 5c                       	popq	%r12
140089b08: 41 5d                       	popq	%r13
140089b0a: 41 5e                       	popq	%r14
140089b0c: 41 5f                       	popq	%r15
140089b0e: c3                          	retq
140089b0f: 90                          	nop
140089b10: 48 85 c0                    	testq	%rax, %rax
140089b13: 0f 85 32 01 00 00           	jne	0x140089c4b <.text+0x88c4b>
140089b19: 31 d2                       	xorl	%edx, %edx
140089b1b: e9 5d fe ff ff              	jmp	0x14008997d <.text+0x8897d>
140089b20: 4c 89 e1                    	movq	%r12, %rcx
140089b23: ff 15 8f c9 0e 00           	callq	*0xec98f(%rip)          # 0x1401764b8
140089b29: 48 8d 0d 26 11 07 00        	leaq	0x71126(%rip), %rcx     # 0x1400fac56
140089b30: e8 4b bb 01 00              	callq	0x1400a5680 <.text+0xa4680>
140089b35: 4c 8d 0d 1a 11 07 00        	leaq	0x7111a(%rip), %r9      # 0x1400fac56
140089b3c: 48 89 44 24 20              	movq	%rax, 0x20(%rsp)
140089b41: eb 85                       	jmp	0x140089ac8 <.text+0x88ac8>
140089b43: 0f 1f 44 00 00              	nopl	(%rax,%rax)
140089b48: ff 15 fa c2 0e 00           	callq	*0xec2fa(%rip)          # 0x140175e48
140089b4e: 48 8b 1d 63 c9 0e 00        	movq	0xec963(%rip), %rbx     # 0x1401764b8
140089b55: 89 c7                       	movl	%eax, %edi
140089b57: 48 89 f1                    	movq	%rsi, %rcx
140089b5a: ff d3                       	callq	*%rbx
140089b5c: 4c 89 e9                    	movq	%r13, %rcx
140089b5f: ff d3                       	callq	*%rbx
140089b61: 4c 89 e1                    	movq	%r12, %rcx
140089b64: ff d3                       	callq	*%rbx
140089b66: 81 ff e2 2e 00 00           	cmpl	$0x2ee2, %edi           # imm = 0x2EE2
140089b6c: 0f 84 5e 01 00 00           	je	0x140089cd0 <.text+0x88cd0>
140089b72: 48 8d 0d 85 10 07 00        	leaq	0x71085(%rip), %rcx     # 0x1400fabfe
140089b79: e8 02 bb 01 00              	callq	0x1400a5680 <.text+0xa4680>
140089b7e: 48 8b bc 24 a0 0b 00 00     	movq	0xba0(%rsp), %rdi
140089b86: 4c 8d 0d 71 10 07 00        	leaq	0x71071(%rip), %r9      # 0x1400fabfe
140089b8d: 4c 8b 47 08                 	movq	0x8(%rdi), %r8
140089b91: 48 89 44 24 20              	movq	%rax, 0x20(%rsp)
140089b96: 48 8b 8c 24 a0 0b 00 00     	movq	0xba0(%rsp), %rcx
140089b9e: 31 d2                       	xorl	%edx, %edx
140089ba0: e8 7b a1 03 00              	callq	0x1400c3d20 <.text+0xc2d20>
140089ba5: 31 db                       	xorl	%ebx, %ebx
140089ba7: 48 8b 8c 24 80 00 00 00     	movq	0x80(%rsp), %rcx
140089baf: 48 8d 84 24 90 00 00 00     	leaq	0x90(%rsp), %rax
140089bb7: 48 39 c1                    	cmpq	%rax, %rcx
140089bba: 0f 84 20 ff ff ff           	je	0x140089ae0 <.text+0x88ae0>
140089bc0: 48 8b 84 24 90 00 00 00     	movq	0x90(%rsp), %rax
140089bc8: 48 8d 54 00 02              	leaq	0x2(%rax,%rax), %rdx
140089bcd: e8 fe 04 04 00              	callq	0x1400ca0d0 <.text+0xc90d0>
140089bd2: e9 09 ff ff ff              	jmp	0x140089ae0 <.text+0x88ae0>
140089bd7: 66 0f 1f 84 00 00 00 00 00  	nopw	(%rax,%rax)
140089be0: 48 39 f9                    	cmpq	%rdi, %rcx
140089be3: 74 11                       	je	0x140089bf6 <.text+0x88bf6>
140089be5: 48 8b 84 24 b0 00 00 00     	movq	0xb0(%rsp), %rax
140089bed: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
140089bf1: e8 da 04 04 00              	callq	0x1400ca0d0 <.text+0xc90d0>
140089bf6: 48 8b 1d bb c8 0e 00        	movq	0xec8bb(%rip), %rbx     # 0x1401764b8
140089bfd: 48 89 f1                    	movq	%rsi, %rcx
140089c00: ff d3                       	callq	*%rbx
140089c02: 4c 89 e9                    	movq	%r13, %rcx
140089c05: ff d3                       	callq	*%rbx
140089c07: 4c 89 e1                    	movq	%r12, %rcx
140089c0a: ff d3                       	callq	*%rbx
140089c0c: bb 01 00 00 00              	movl	$0x1, %ebx
140089c11: eb 94                       	jmp	0x140089ba7 <.text+0x88ba7>
140089c13: 0f 1f 44 00 00              	nopl	(%rax,%rax)
140089c18: 48 8d 94 24 a0 00 00 00     	leaq	0xa0(%rsp), %rdx
140089c20: 45 31 c0                    	xorl	%r8d, %r8d
140089c23: 48 8d 8c 24 80 00 00 00     	leaq	0x80(%rsp), %rcx
140089c2b: e8 b0 cb 03 00              	callq	0x1400c67e0 <.text+0xc57e0>
140089c30: 48 89 84 24 80 00 00 00     	movq	%rax, 0x80(%rsp)
140089c38: 48 89 c1                    	movq	%rax, %rcx
140089c3b: 48 8b 84 24 a0 00 00 00     	movq	0xa0(%rsp), %rax
140089c43: 48 89 84 24 90 00 00 00     	movq	%rax, 0x90(%rsp)
140089c4b: 48 8d 15 2e 10 07 00        	leaq	0x7102e(%rip), %rdx     # 0x1400fac80
140089c52: 4d 89 f0                    	movq	%r14, %r8
140089c55: e8 8e bb 01 00              	callq	0x1400a57e8 <.text+0xa47e8>
140089c5a: 48 8b 84 24 a0 00 00 00     	movq	0xa0(%rsp), %rax
140089c62: 48 8b 8c 24 80 00 00 00     	movq	0x80(%rsp), %rcx
140089c6a: 48 8d 14 00                 	leaq	(%rax,%rax), %rdx
140089c6e: e9 0a fd ff ff              	jmp	0x14008997d <.text+0x8897d>
140089c73: 0f 1f 44 00 00              	nopl	(%rax,%rax)
140089c78: 48 8d 0d b9 0f 07 00        	leaq	0x70fb9(%rip), %rcx     # 0x1400fac38
140089c7f: e8 fc b9 01 00              	callq	0x1400a5680 <.text+0xa4680>
140089c84: 4c 8d 0d ad 0f 07 00        	leaq	0x70fad(%rip), %r9      # 0x1400fac38
140089c8b: 48 89 44 24 20              	movq	%rax, 0x20(%rsp)
140089c90: e9 33 fe ff ff              	jmp	0x140089ac8 <.text+0x88ac8>
140089c95: 0f 1f 00                    	nopl	(%rax)
140089c98: 48 8b 1d 19 c8 0e 00        	movq	0xec819(%rip), %rbx     # 0x1401764b8
140089c9f: 4c 89 e9                    	movq	%r13, %rcx
140089ca2: ff d3                       	callq	*%rbx
140089ca4: 4c 89 e1                    	movq	%r12, %rcx
140089ca7: ff d3                       	callq	*%rbx
140089ca9: 48 8d 0d c0 0f 07 00        	leaq	0x70fc0(%rip), %rcx     # 0x1400fac70
140089cb0: e8 cb b9 01 00              	callq	0x1400a5680 <.text+0xa4680>
140089cb5: 4c 8d 0d b4 0f 07 00        	leaq	0x70fb4(%rip), %r9      # 0x1400fac70
140089cbc: 48 89 44 24 20              	movq	%rax, 0x20(%rsp)
140089cc1: e9 02 fe ff ff              	jmp	0x140089ac8 <.text+0x88ac8>
140089cc6: 66 2e 0f 1f 84 00 00 00 00 00       	nopw	%cs:(%rax,%rax)
140089cd0: 48 8d 0d 31 0f 07 00        	leaq	0x70f31(%rip), %rcx     # 0x1400fac08
140089cd7: e8 a4 b9 01 00              	callq	0x1400a5680 <.text+0xa4680>
140089cdc: 48 8b bc 24 a0 0b 00 00     	movq	0xba0(%rsp), %rdi
140089ce4: 4c 8d 0d 1d 0f 07 00        	leaq	0x70f1d(%rip), %r9      # 0x1400fac08
140089ceb: 4c 8b 47 08                 	movq	0x8(%rdi), %r8
140089cef: e9 9d fe ff ff              	jmp	0x140089b91 <.text+0x88b91>
140089cf4: 48 89 c3                    	movq	%rax, %rbx
140089cf7: e9 1b e4 04 00              	jmp	0x1400d8117 <.text+0xd7117>
140089cfc: e9 35 e4 04 00              	jmp	0x1400d8136 <.text+0xd7136>
140089d01: e9 01 e4 04 00              	jmp	0x1400d8107 <.text+0xd7107>
140089d06: 48 89 c3                    	movq	%rax, %rbx
140089d09: e9 16 e4 04 00              	jmp	0x1400d8124 <.text+0xd7124>
140089d0e: 66 90                       	nop
140089d10: 53                          	pushq	%rbx
140089d11: 48 83 ec 40                 	subq	$0x40, %rsp
140089d15: 4c 8b 05 4c a6 0e 00        	movq	0xea64c(%rip), %r8      # 0x140174368
140089d1c: 48 89 cb                    	movq	%rcx, %rbx
140089d1f: 4d 85 c0                    	testq	%r8, %r8
140089d22: 0f 84 90 00 00 00           	je	0x140089db8 <.text+0x88db8>
140089d28: 48 8d 49 10                 	leaq	0x10(%rcx), %rcx
140089d2c: 4c 89 44 24 38              	movq	%r8, 0x38(%rsp)
140089d31: 48 89 0b                    	movq	%rcx, (%rbx)
140089d34: 4c 8b 0d 25 a6 0e 00        	movq	0xea625(%rip), %r9      # 0x140174360
140089d3b: 49 83 f8 07                 	cmpq	$0x7, %r8
140089d3f: 77 27                       	ja	0x140089d68 <.text+0x88d68>
140089d41: 49 83 c0 01                 	addq	$0x1, %r8
140089d45: 4d 01 c0                    	addq	%r8, %r8
140089d48: 4c 89 ca                    	movq	%r9, %rdx
140089d4b: e8 98 ba 01 00              	callq	0x1400a57e8 <.text+0xa47e8>
140089d50: 48 8b 44 24 38              	movq	0x38(%rsp), %rax
140089d55: 48 89 43 08                 	movq	%rax, 0x8(%rbx)
140089d59: 48 89 d8                    	movq	%rbx, %rax
140089d5c: 48 83 c4 40                 	addq	$0x40, %rsp
140089d60: 5b                          	popq	%rbx
140089d61: c3                          	retq
140089d62: 66 0f 1f 44 00 00           	nopw	(%rax,%rax)
140089d68: 45 31 c0                    	xorl	%r8d, %r8d
140089d6b: 48 89 d9                    	movq	%rbx, %rcx
140089d6e: 48 8d 54 24 38              	leaq	0x38(%rsp), %rdx
140089d73: 4c 89 4c 24 28              	movq	%r9, 0x28(%rsp)
140089d78: e8 63 ca 03 00              	callq	0x1400c67e0 <.text+0xc57e0>
140089d7d: 4c 8b 44 24 38              	movq	0x38(%rsp), %r8
140089d82: 4c 8b 4c 24 28              	movq	0x28(%rsp), %r9
140089d87: 48 89 03                    	movq	%rax, (%rbx)
140089d8a: 48 89 c1                    	movq	%rax, %rcx
140089d8d: 4d 85 c0                    	testq	%r8, %r8
140089d90: 4c 89 43 10                 	movq	%r8, 0x10(%rbx)
140089d94: 74 12                       	je	0x140089da8 <.text+0x88da8>
140089d96: 48 c7 c0 ff ff ff ff        	movq	$-0x1, %rax
140089d9d: 49 83 c0 01                 	addq	$0x1, %r8
140089da1: 74 b2                       	je	0x140089d55 <.text+0x88d55>
140089da3: eb a0                       	jmp	0x140089d45 <.text+0x88d45>
140089da5: 0f 1f 00                    	nopl	(%rax)
140089da8: 41 0f b7 01                 	movzwl	(%r9), %eax
140089dac: 66 89 01                    	movw	%ax, (%rcx)
140089daf: 31 c0                       	xorl	%eax, %eax
140089db1: eb a2                       	jmp	0x140089d55 <.text+0x88d55>
140089db3: 0f 1f 44 00 00              	nopl	(%rax,%rax)
140089db8: e8 73 63 00 00              	callq	0x140090130 <.text+0x8f130>
140089dbd: 48 89 d8                    	movq	%rbx, %rax
140089dc0: 48 83 c4 40                 	addq	$0x40, %rsp
140089dc4: 5b                          	popq	%rbx
140089dc5: c3                          	retq
140089dc6: 66 2e 0f 1f 84 00 00 00 00 00       	nopw	%cs:(%rax,%rax)
140089dd0: 48 83 ec 28                 	subq	$0x28, %rsp
140089dd4: 48 89 ca                    	movq	%rcx, %rdx
140089dd7: 48 8d 0d 82 a5 0e 00        	leaq	0xea582(%rip), %rcx     # 0x140174360
140089dde: e8 dd c8 03 00              	callq	0x1400c66c0 <.text+0xc56c0>
140089de3: 48 8b 15 7e a5 0e 00        	movq	0xea57e(%rip), %rdx     # 0x140174368
140089dea: 48 85 d2                    	testq	%rdx, %rdx
140089ded: 74 3a                       	je	0x140089e29 <.text+0x88e29>
140089def: 90                          	nop
140089df0: 48 8b 05 69 a5 0e 00        	movq	0xea569(%rip), %rax     # 0x140174360
140089df7: 48 83 ea 01                 	subq	$0x1, %rdx
140089dfb: 0f b7 04 50                 	movzwl	(%rax,%rdx,2), %eax
140089dff: 66 83 f8 5c                 	cmpw	$0x5c, %ax
140089e03: 74 06                       	je	0x140089e0b <.text+0x88e0b>
140089e05: 66 83 f8 2f                 	cmpw	$0x2f, %ax
140089e09: 75 1e                       	jne	0x140089e29 <.text+0x88e29>
140089e0b: 41 b8 01 00 00 00           	movl	$0x1, %r8d
140089e11: 48 8d 0d 48 a5 0e 00        	leaq	0xea548(%rip), %rcx     # 0x140174360
140089e18: e8 43 c7 03 00              	callq	0x1400c6560 <.text+0xc5560>
140089e1d: 48 8b 15 44 a5 0e 00        	movq	0xea544(%rip), %rdx     # 0x140174368
140089e24: 48 85 d2                    	testq	%rdx, %rdx
140089e27: 75 c7                       	jne	0x140089df0 <.text+0x88df0>
140089e29: 48 83 c4 28                 	addq	$0x28, %rsp
140089e2d: c3                          	retq
140089e2e: 66 90                       	nop
140089e30: 41 56                       	pushq	%r14
140089e32: 41 55                       	pushq	%r13
140089e34: 41 54                       	pushq	%r12
140089e36: 55                          	pushq	%rbp
140089e37: 57                          	pushq	%rdi
140089e38: 56                          	pushq	%rsi
140089e39: 53                          	pushq	%rbx
140089e3a: 48 81 ec e0 00 00 00        	subq	$0xe0, %rsp
140089e41: 48 8d 4c 24 40              	leaq	0x40(%rsp), %rcx
140089e46: e8 c5 fe ff ff              	callq	0x140089d10 <.text+0x88d10>
140089e4b: 48 83 7c 24 48 00           	cmpq	$0x0, 0x48(%rsp)
140089e51: 75 3d                       	jne	0x140089e90 <.text+0x88e90>
140089e53: 48 8b 4c 24 40              	movq	0x40(%rsp), %rcx
140089e58: 48 8d 44 24 50              	leaq	0x50(%rsp), %rax
140089e5d: 31 db                       	xorl	%ebx, %ebx
140089e5f: 48 39 c1                    	cmpq	%rax, %rcx
140089e62: 74 0f                       	je	0x140089e73 <.text+0x88e73>
140089e64: 48 8b 44 24 50              	movq	0x50(%rsp), %rax
140089e69: 48 8d 54 00 02              	leaq	0x2(%rax,%rax), %rdx
140089e6e: e8 5d 02 04 00              	callq	0x1400ca0d0 <.text+0xc90d0>
140089e73: 89 d8                       	movl	%ebx, %eax
140089e75: 48 81 c4 e0 00 00 00        	addq	$0xe0, %rsp
140089e7c: 5b                          	popq	%rbx
140089e7d: 5e                          	popq	%rsi
140089e7e: 5f                          	popq	%rdi
140089e7f: 5d                          	popq	%rbp
140089e80: 41 5c                       	popq	%r12
140089e82: 41 5d                       	popq	%r13
140089e84: 41 5e                       	popq	%r14
140089e86: c3                          	retq
140089e87: 66 0f 1f 84 00 00 00 00 00  	nopw	(%rax,%rax)
140089e90: 48 8d 15 2b 0e 07 00        	leaq	0x70e2b(%rip), %rdx     # 0x1400facc2
140089e97: 48 8d 74 24 70              	leaq	0x70(%rsp), %rsi
140089e9c: 45 31 d2                    	xorl	%r10d, %r10d
140089e9f: c6 44 24 70 00              	movb	$0x0, 0x70(%rsp)
140089ea4: 4c 8d 42 02                 	leaq	0x2(%rdx), %r8
140089ea8: 48 8d bc 24 90 00 00 00     	leaq	0x90(%rsp), %rdi
140089eb0: 48 89 74 24 60              	movq	%rsi, 0x60(%rsp)
140089eb5: 48 8d ac 24 d0 00 00 00     	leaq	0xd0(%rsp), %rbp
140089ebd: 48 8d 8c 24 c0 00 00 00     	leaq	0xc0(%rsp), %rcx
140089ec5: 48 c7 44 24 68 00 00 00 00  	movq	$0x0, 0x68(%rsp)
140089ece: 48 89 bc 24 80 00 00 00     	movq	%rdi, 0x80(%rsp)
140089ed6: 4c 8d 6c 24 60              	leaq	0x60(%rsp), %r13
140089edb: 4c 8d b4 24 80 00 00 00     	leaq	0x80(%rsp), %r14
140089ee3: 66 44 89 94 24 90 00 00 00  	movw	%r10w, 0x90(%rsp)
140089eec: 48 c7 84 24 88 00 00 00 00 00 00 00 	movq	$0x0, 0x88(%rsp)
140089ef8: 48 89 ac 24 c0 00 00 00     	movq	%rbp, 0xc0(%rsp)
140089f00: e8 1b ec ff ff              	callq	0x140088b20 <.text+0x87b20>
140089f05: 4c 8d a4 24 b0 00 00 00     	leaq	0xb0(%rsp), %r12
140089f0d: 48 8d 0d b2 0d 07 00        	leaq	0x70db2(%rip), %rcx     # 0x1400facc6
140089f14: 4c 89 a4 24 a0 00 00 00     	movq	%r12, 0xa0(%rsp)
140089f1c: e8 5f b7 01 00              	callq	0x1400a5680 <.text+0xa4680>
140089f21: 48 8d 1c 00                 	leaq	(%rax,%rax), %rbx
140089f25: 48 89 da                    	movq	%rbx, %rdx
140089f28: 48 d1 fa                    	sarq	%rdx
140089f2b: 48 89 54 24 38              	movq	%rdx, 0x38(%rsp)
140089f30: 48 83 fb 0e                 	cmpq	$0xe, %rbx
140089f34: 0f 87 16 01 00 00           	ja	0x14008a050 <.text+0x89050>
140089f3a: 48 8b 8c 24 a0 00 00 00     	movq	0xa0(%rsp), %rcx
140089f42: 48 83 fa 01                 	cmpq	$0x1, %rdx
140089f46: 0f 85 f4 00 00 00           	jne	0x14008a040 <.text+0x89040>
140089f4c: 41 b9 2f 00 00 00           	movl	$0x2f, %r9d
140089f52: b8 02 00 00 00              	movl	$0x2, %eax
140089f57: 66 44 89 09                 	movw	%r9w, (%rcx)
140089f5b: 45 31 c0                    	xorl	%r8d, %r8d
140089f5e: 4c 8d b4 24 80 00 00 00     	leaq	0x80(%rsp), %r14
140089f66: 4c 8d 6c 24 60              	leaq	0x60(%rsp), %r13
140089f6b: 66 44 89 04 01              	movw	%r8w, (%rcx,%rax)
140089f70: 4d 89 e9                    	movq	%r13, %r9
140089f73: 48 8d 4c 24 40              	leaq	0x40(%rsp), %rcx
140089f78: 4c 8d 84 24 c0 00 00 00     	leaq	0xc0(%rsp), %r8
140089f80: 4c 89 74 24 20              	movq	%r14, 0x20(%rsp)
140089f85: 48 89 94 24 a8 00 00 00     	movq	%rdx, 0xa8(%rsp)
140089f8d: 48 8d 94 24 a0 00 00 00     	leaq	0xa0(%rsp), %rdx
140089f95: e8 56 f7 ff ff              	callq	0x1400896f0 <.text+0x886f0>
140089f9a: 48 8b 8c 24 a0 00 00 00     	movq	0xa0(%rsp), %rcx
140089fa2: 89 c3                       	movl	%eax, %ebx
140089fa4: 4c 39 e1                    	cmpq	%r12, %rcx
140089fa7: 74 12                       	je	0x140089fbb <.text+0x88fbb>
140089fa9: 48 8b 84 24 b0 00 00 00     	movq	0xb0(%rsp), %rax
140089fb1: 48 8d 54 00 02              	leaq	0x2(%rax,%rax), %rdx
140089fb6: e8 15 01 04 00              	callq	0x1400ca0d0 <.text+0xc90d0>
140089fbb: 48 8b 8c 24 c0 00 00 00     	movq	0xc0(%rsp), %rcx
140089fc3: 48 39 e9                    	cmpq	%rbp, %rcx
140089fc6: 74 11                       	je	0x140089fd9 <.text+0x88fd9>
140089fc8: 48 8b 84 24 d0 00 00 00     	movq	0xd0(%rsp), %rax
140089fd0: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
140089fd4: e8 f7 00 04 00              	callq	0x1400ca0d0 <.text+0xc90d0>
140089fd9: 84 db                       	testb	%bl, %bl
140089fdb: 74 0a                       	je	0x140089fe7 <.text+0x88fe7>
140089fdd: 4c 89 e9                    	movq	%r13, %rcx
140089fe0: e8 9b f6 ff ff              	callq	0x140089680 <.text+0x88680>
140089fe5: 89 c3                       	movl	%eax, %ebx
140089fe7: 48 8b 8c 24 80 00 00 00     	movq	0x80(%rsp), %rcx
140089fef: 48 39 f9                    	cmpq	%rdi, %rcx
140089ff2: 74 12                       	je	0x14008a006 <.text+0x89006>
140089ff4: 48 8b 84 24 90 00 00 00     	movq	0x90(%rsp), %rax
140089ffc: 48 8d 54 00 02              	leaq	0x2(%rax,%rax), %rdx
14008a001: e8 ca 00 04 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008a006: 48 8b 4c 24 60              	movq	0x60(%rsp), %rcx
14008a00b: 48 39 f1                    	cmpq	%rsi, %rcx
14008a00e: 74 0e                       	je	0x14008a01e <.text+0x8901e>
14008a010: 48 8b 44 24 70              	movq	0x70(%rsp), %rax
14008a015: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008a019: e8 b2 00 04 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008a01e: 48 8b 4c 24 40              	movq	0x40(%rsp), %rcx
14008a023: 48 8d 44 24 50              	leaq	0x50(%rsp), %rax
14008a028: 48 39 c1                    	cmpq	%rax, %rcx
14008a02b: 0f 85 33 fe ff ff           	jne	0x140089e64 <.text+0x88e64>
14008a031: e9 3d fe ff ff              	jmp	0x140089e73 <.text+0x88e73>
14008a036: 66 2e 0f 1f 84 00 00 00 00 00       	nopw	%cs:(%rax,%rax)
14008a040: 48 85 d2                    	testq	%rdx, %rdx
14008a043: 75 38                       	jne	0x14008a07d <.text+0x8907d>
14008a045: 31 c0                       	xorl	%eax, %eax
14008a047: e9 0f ff ff ff              	jmp	0x140089f5b <.text+0x88f5b>
14008a04c: 0f 1f 40 00                 	nopl	(%rax)
14008a050: 48 8d 54 24 38              	leaq	0x38(%rsp), %rdx
14008a055: 45 31 c0                    	xorl	%r8d, %r8d
14008a058: 48 8d 8c 24 a0 00 00 00     	leaq	0xa0(%rsp), %rcx
14008a060: e8 7b c7 03 00              	callq	0x1400c67e0 <.text+0xc57e0>
14008a065: 48 89 84 24 a0 00 00 00     	movq	%rax, 0xa0(%rsp)
14008a06d: 48 89 c1                    	movq	%rax, %rcx
14008a070: 48 8b 44 24 38              	movq	0x38(%rsp), %rax
14008a075: 48 89 84 24 b0 00 00 00     	movq	%rax, 0xb0(%rsp)
14008a07d: 48 8d 15 42 0c 07 00        	leaq	0x70c42(%rip), %rdx     # 0x1400facc6
14008a084: 49 89 d8                    	movq	%rbx, %r8
14008a087: e8 5c b7 01 00              	callq	0x1400a57e8 <.text+0xa47e8>
14008a08c: 48 8b 54 24 38              	movq	0x38(%rsp), %rdx
14008a091: 48 8b 8c 24 a0 00 00 00     	movq	0xa0(%rsp), %rcx
14008a099: 48 8d 04 12                 	leaq	(%rdx,%rdx), %rax
14008a09d: e9 b9 fe ff ff              	jmp	0x140089f5b <.text+0x88f5b>
14008a0a2: 48 89 c3                    	movq	%rax, %rbx
14008a0a5: e9 ae e0 04 00              	jmp	0x1400d8158 <.text+0xd7158>
14008a0aa: 48 89 c3                    	movq	%rax, %rbx
14008a0ad: e9 c0 e0 04 00              	jmp	0x1400d8172 <.text+0xd7172>
14008a0b2: 48 89 c3                    	movq	%rax, %rbx
14008a0b5: 4c 8d 6c 24 60              	leaq	0x60(%rsp), %r13
14008a0ba: 4c 8d b4 24 80 00 00 00     	leaq	0x80(%rsp), %r14
14008a0c2: e9 9e e0 04 00              	jmp	0x1400d8165 <.text+0xd7165>
14008a0c7: 66 0f 1f 84 00 00 00 00 00  	nopw	(%rax,%rax)
14008a0d0: 41 57                       	pushq	%r15
14008a0d2: 41 56                       	pushq	%r14
14008a0d4: 41 55                       	pushq	%r13
14008a0d6: 41 54                       	pushq	%r12
14008a0d8: 55                          	pushq	%rbp
14008a0d9: 57                          	pushq	%rdi
14008a0da: 56                          	pushq	%rsi
14008a0db: 53                          	pushq	%rbx
14008a0dc: 48 81 ec 38 03 00 00        	subq	$0x338, %rsp            # imm = 0x338
14008a0e3: 48 89 cb                    	movq	%rcx, %rbx
14008a0e6: 48 8d 4c 24 60              	leaq	0x60(%rsp), %rcx
14008a0eb: 48 89 d7                    	movq	%rdx, %rdi
14008a0ee: 4c 89 c5                    	movq	%r8, %rbp
14008a0f1: 4d 89 cc                    	movq	%r9, %r12
14008a0f4: e8 17 fc ff ff              	callq	0x140089d10 <.text+0x88d10>
14008a0f9: 48 83 7c 24 68 00           	cmpq	$0x0, 0x68(%rsp)
14008a0ff: 0f 84 3b 0b 00 00           	je	0x14008ac40 <.text+0x89c40>
14008a105: 49 83 fc 09                 	cmpq	$0x9, %r12
14008a109: 0f 86 3c 10 00 00           	jbe	0x14008b14b <.text+0x8a14b>
14008a10f: 4c 89 e1                    	movq	%r12, %rcx
14008a112: 41 be 01 00 00 00           	movl	$0x1, %r14d
14008a118: 49 b8 4b 59 86 38 d6 c5 6d 34       	movabsq	$0x346dc5d63886594b, %r8 # imm = 0x346DC5D63886594B
14008a122: eb 3c                       	jmp	0x14008a160 <.text+0x89160>
14008a124: 0f 1f 40 00                 	nopl	(%rax)
14008a128: 48 81 f9 e7 03 00 00        	cmpq	$0x3e7, %rcx            # imm = 0x3E7
14008a12f: 0f 86 bb 0b 00 00           	jbe	0x14008acf0 <.text+0x89cf0>
14008a135: 48 81 f9 0f 27 00 00        	cmpq	$0x270f, %rcx           # imm = 0x270F
14008a13c: 0f 86 be 0b 00 00           	jbe	0x14008ad00 <.text+0x89d00>
14008a142: 48 89 c8                    	movq	%rcx, %rax
14008a145: 41 83 c6 04                 	addl	$0x4, %r14d
14008a149: 49 f7 e0                    	mulq	%r8
14008a14c: 48 81 f9 9f 86 01 00        	cmpq	$0x1869f, %rcx          # imm = 0x1869F
14008a153: 0f 86 b7 0b 00 00           	jbe	0x14008ad10 <.text+0x89d10>
14008a159: 48 c1 ea 0b                 	shrq	$0xb, %rdx
14008a15d: 48 89 d1                    	movq	%rdx, %rcx
14008a160: 48 83 f9 63                 	cmpq	$0x63, %rcx
14008a164: 77 c2                       	ja	0x14008a128 <.text+0x89128>
14008a166: 45 8d 6e 01                 	leal	0x1(%r14), %r13d
14008a16a: 4d 89 ee                    	movq	%r13, %r14
14008a16d: 48 8d b4 24 50 02 00 00     	leaq	0x250(%rsp), %rsi
14008a175: 4c 89 ea                    	movq	%r13, %rdx
14008a178: 48 8d 8c 24 40 02 00 00     	leaq	0x240(%rsp), %rcx
14008a180: c6 84 24 50 02 00 00 00     	movb	$0x0, 0x250(%rsp)
14008a188: 48 89 b4 24 40 02 00 00     	movq	%rsi, 0x240(%rsp)
14008a190: 48 c7 84 24 48 02 00 00 00 00 00 00 	movq	$0x0, 0x248(%rsp)
14008a19c: e8 ff 85 03 00              	callq	0x1400c27a0 <.text+0xc17a0>
14008a1a1: 66 0f 6f 05 f7 0c 07 00     	movdqa	0x70cf7(%rip), %xmm0    # 0x1400faea0
14008a1a9: 4c 8b 8c 24 40 02 00 00     	movq	0x240(%rsp), %r9
14008a1b1: 0f 29 84 24 60 02 00 00     	movaps	%xmm0, 0x260(%rsp)
14008a1b9: 66 0f 6f 05 ef 0c 07 00     	movdqa	0x70cef(%rip), %xmm0    # 0x1400faeb0
14008a1c1: 0f 29 84 24 70 02 00 00     	movaps	%xmm0, 0x270(%rsp)
14008a1c9: 66 0f 6f 05 ef 0c 07 00     	movdqa	0x70cef(%rip), %xmm0    # 0x1400faec0
14008a1d1: 0f 29 84 24 80 02 00 00     	movaps	%xmm0, 0x280(%rsp)
14008a1d9: 66 0f 6f 05 ef 0c 07 00     	movdqa	0x70cef(%rip), %xmm0    # 0x1400faed0
14008a1e1: 0f 29 84 24 90 02 00 00     	movaps	%xmm0, 0x290(%rsp)
14008a1e9: 66 0f 6f 05 ef 0c 07 00     	movdqa	0x70cef(%rip), %xmm0    # 0x1400faee0
14008a1f1: 0f 29 84 24 a0 02 00 00     	movaps	%xmm0, 0x2a0(%rsp)
14008a1f9: 66 0f 6f 05 ef 0c 07 00     	movdqa	0x70cef(%rip), %xmm0    # 0x1400faef0
14008a201: 0f 29 84 24 b0 02 00 00     	movaps	%xmm0, 0x2b0(%rsp)
14008a209: 66 0f 6f 05 ef 0c 07 00     	movdqa	0x70cef(%rip), %xmm0    # 0x1400faf00
14008a211: 0f 29 84 24 c0 02 00 00     	movaps	%xmm0, 0x2c0(%rsp)
14008a219: 66 0f 6f 05 ef 0c 07 00     	movdqa	0x70cef(%rip), %xmm0    # 0x1400faf10
14008a221: 0f 29 84 24 d0 02 00 00     	movaps	%xmm0, 0x2d0(%rsp)
14008a229: 66 0f 6f 05 ef 0c 07 00     	movdqa	0x70cef(%rip), %xmm0    # 0x1400faf20
14008a231: 0f 29 84 24 e0 02 00 00     	movaps	%xmm0, 0x2e0(%rsp)
14008a239: 66 0f 6f 05 ef 0c 07 00     	movdqa	0x70cef(%rip), %xmm0    # 0x1400faf30
14008a241: 0f 29 84 24 f0 02 00 00     	movaps	%xmm0, 0x2f0(%rsp)
14008a249: 66 0f 6f 05 ef 0c 07 00     	movdqa	0x70cef(%rip), %xmm0    # 0x1400faf40
14008a251: 0f 29 84 24 00 03 00 00     	movaps	%xmm0, 0x300(%rsp)
14008a259: 66 0f 6f 05 ef 0c 07 00     	movdqa	0x70cef(%rip), %xmm0    # 0x1400faf50
14008a261: 0f 29 84 24 10 03 00 00     	movaps	%xmm0, 0x310(%rsp)
14008a269: 66 0f 6f 05 ef 0c 07 00     	movdqa	0x70cef(%rip), %xmm0    # 0x1400faf60
14008a271: 0f 11 84 24 19 03 00 00     	movups	%xmm0, 0x319(%rsp)
14008a279: 49 83 fc 63                 	cmpq	$0x63, %r12
14008a27d: 76 5f                       	jbe	0x14008a2de <.text+0x892de>
14008a27f: 49 b8 c3 f5 28 5c 8f c2 f5 28       	movabsq	$0x28f5c28f5c28f5c3, %r8 # imm = 0x28F5C28F5C28F5C3
14008a289: 41 83 ee 01                 	subl	$0x1, %r14d
14008a28d: 0f 1f 00                    	nopl	(%rax)
14008a290: 4c 89 e2                    	movq	%r12, %rdx
14008a293: 48 c1 ea 02                 	shrq	$0x2, %rdx
14008a297: 48 89 d0                    	movq	%rdx, %rax
14008a29a: 49 f7 e0                    	mulq	%r8
14008a29d: 4c 89 e0                    	movq	%r12, %rax
14008a2a0: 48 c1 ea 02                 	shrq	$0x2, %rdx
14008a2a4: 48 6b ca 64                 	imulq	$0x64, %rdx, %rcx
14008a2a8: 48 29 c8                    	subq	%rcx, %rax
14008a2ab: 4c 89 e1                    	movq	%r12, %rcx
14008a2ae: 49 89 d4                    	movq	%rdx, %r12
14008a2b1: 44 89 f2                    	movl	%r14d, %edx
14008a2b4: 44 0f b6 94 44 61 02 00 00  	movzbl	0x261(%rsp,%rax,2), %r10d
14008a2bd: 45 88 14 11                 	movb	%r10b, (%r9,%rdx)
14008a2c1: 0f b6 84 44 60 02 00 00     	movzbl	0x260(%rsp,%rax,2), %eax
14008a2c9: 41 8d 56 ff                 	leal	-0x1(%r14), %edx
14008a2cd: 41 83 ee 02                 	subl	$0x2, %r14d
14008a2d1: 41 88 04 11                 	movb	%al, (%r9,%rdx)
14008a2d5: 48 81 f9 0f 27 00 00        	cmpq	$0x270f, %rcx           # imm = 0x270F
14008a2dc: 77 b2                       	ja	0x14008a290 <.text+0x89290>
14008a2de: 41 8d 44 24 30              	leal	0x30(%r12), %eax
14008a2e3: 49 83 fc 09                 	cmpq	$0x9, %r12
14008a2e7: 76 16                       	jbe	0x14008a2ff <.text+0x892ff>
14008a2e9: 42 0f b6 84 64 61 02 00 00  	movzbl	0x261(%rsp,%r12,2), %eax
14008a2f2: 41 88 41 01                 	movb	%al, 0x1(%r9)
14008a2f6: 42 0f b6 84 64 60 02 00 00  	movzbl	0x260(%rsp,%r12,2), %eax
14008a2ff: 41 88 01                    	movb	%al, (%r9)
14008a302: 48 8b 84 24 40 02 00 00     	movq	0x240(%rsp), %rax
14008a30a: 48 89 ea                    	movq	%rbp, %rdx
14008a30d: 48 8d 8c 24 c0 01 00 00     	leaq	0x1c0(%rsp), %rcx
14008a315: 42 c6 04 28 00              	movb	$0x0, (%rax,%r13)
14008a31a: 4c 89 ac 24 48 02 00 00     	movq	%r13, 0x248(%rsp)
14008a322: e8 29 eb ff ff              	callq	0x140088e50 <.text+0x87e50>
14008a327: 48 8d 94 24 c0 01 00 00     	leaq	0x1c0(%rsp), %rdx
14008a32f: 48 8d 8c 24 e0 01 00 00     	leaq	0x1e0(%rsp), %rcx
14008a337: e8 d4 e8 ff ff              	callq	0x140088c10 <.text+0x87c10>
14008a33c: 48 89 fa                    	movq	%rdi, %rdx
14008a33f: 48 8d 8c 24 40 01 00 00     	leaq	0x140(%rsp), %rcx
14008a347: e8 04 eb ff ff              	callq	0x140088e50 <.text+0x87e50>
14008a34c: 48 8d 94 24 40 01 00 00     	leaq	0x140(%rsp), %rdx
14008a354: 48 8d 8c 24 60 01 00 00     	leaq	0x160(%rsp), %rcx
14008a35c: e8 af e8 ff ff              	callq	0x140088c10 <.text+0x87c10>
14008a361: 48 89 da                    	movq	%rbx, %rdx
14008a364: 48 8d 8c 24 c0 00 00 00     	leaq	0xc0(%rsp), %rcx
14008a36c: e8 df ea ff ff              	callq	0x140088e50 <.text+0x87e50>
14008a371: 48 8d 94 24 c0 00 00 00     	leaq	0xc0(%rsp), %rdx
14008a379: 48 8d 8c 24 e0 00 00 00     	leaq	0xe0(%rsp), %rcx
14008a381: e8 8a e8 ff ff              	callq	0x140088c10 <.text+0x87c10>
14008a386: 48 8d 15 57 09 07 00        	leaq	0x70957(%rip), %rdx     # 0x1400face4
14008a38d: 4c 8d b4 24 b0 00 00 00     	leaq	0xb0(%rsp), %r14
14008a395: 4c 8d 42 0a                 	leaq	0xa(%rdx), %r8
14008a399: 48 8d 8c 24 a0 00 00 00     	leaq	0xa0(%rsp), %rcx
14008a3a1: 4c 89 b4 24 a0 00 00 00     	movq	%r14, 0xa0(%rsp)
14008a3a9: e8 72 e7 ff ff              	callq	0x140088b20 <.text+0x87b20>
14008a3ae: 48 8b 84 24 a8 00 00 00     	movq	0xa8(%rsp), %rax
14008a3b6: 4c 8b 84 24 e8 00 00 00     	movq	0xe8(%rsp), %r8
14008a3be: 4c 8b 8c 24 a0 00 00 00     	movq	0xa0(%rsp), %r9
14008a3c6: 48 8b 8c 24 e0 00 00 00     	movq	0xe0(%rsp), %rcx
14008a3ce: 4a 8d 14 00                 	leaq	(%rax,%r8), %rdx
14008a3d2: 4d 39 f1                    	cmpq	%r14, %r9
14008a3d5: 0f 84 05 0c 00 00           	je	0x14008afe0 <.text+0x89fe0>
14008a3db: 48 39 94 24 b0 00 00 00     	cmpq	%rdx, 0xb0(%rsp)
14008a3e3: 73 1b                       	jae	0x14008a400 <.text+0x89400>
14008a3e5: 4c 8d 94 24 f0 00 00 00     	leaq	0xf0(%rsp), %r10
14008a3ed: 4c 39 d1                    	cmpq	%r10, %rcx
14008a3f0: 74 0e                       	je	0x14008a400 <.text+0x89400>
14008a3f2: 48 39 94 24 f0 00 00 00     	cmpq	%rdx, 0xf0(%rsp)
14008a3fa: 0f 83 a0 09 00 00           	jae	0x14008ada0 <.text+0x89da0>
