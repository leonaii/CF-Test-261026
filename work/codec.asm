
CF上号器_v1.17.exe:	file format coff-x86-64

Disassembly of section .text:

0000000140001000 <.text>:
140088e00: 44 24 47                    	andb	$0x47, %al
140088e03: 48 89 54 24 38              	movq	%rdx, 0x38(%rsp)
140088e08: e8 03 9d 03 00              	callq	0x1400c2b10 <.text+0xc1b10>
140088e0d: 48 8b 0e                    	movq	(%rsi), %rcx
140088e10: 0f b6 44 24 47              	movzbl	0x47(%rsp), %eax
140088e15: 48 8b 54 24 38              	movq	0x38(%rsp), %rdx
140088e1a: 4c 8b 54 24 48              	movq	0x48(%rsp), %r10
140088e1f: 88 04 11                    	movb	%al, (%rcx,%rdx)
140088e22: 48 8b 06                    	movq	(%rsi), %rax
140088e25: c6 44 10 01 00              	movb	$0x0, 0x1(%rax,%rdx)
140088e2a: 4c 89 56 08                 	movq	%r10, 0x8(%rsi)
140088e2e: e9 6a fe ff ff              	jmp	0x140088c9d <.text+0x87c9d>
140088e33: 0f 1f 44 00 00              	nopl	(%rax,%rax)
140088e38: 41 b8 0f 00 00 00           	movl	$0xf, %r8d
140088e3e: e9 11 ff ff ff              	jmp	0x140088d54 <.text+0x87d54>
140088e43: e9 e0 f1 04 00              	jmp	0x1400d8028 <.text+0xd7028>
140088e48: 0f 1f 84 00 00 00 00 00     	nopl	(%rax,%rax)
140088e50: 56                          	pushq	%rsi
140088e51: 53                          	pushq	%rbx
140088e52: 48 83 ec 78                 	subq	$0x78, %rsp
140088e56: 4c 8b 4a 08                 	movq	0x8(%rdx), %r9
140088e5a: 49 89 ca                    	movq	%rcx, %r10
140088e5d: 4d 85 c9                    	testq	%r9, %r9
140088e60: 75 1e                       	jne	0x140088e80 <.text+0x87e80>
140088e62: 48 8d 41 10                 	leaq	0x10(%rcx), %rax
140088e66: c6 41 10 00                 	movb	$0x0, 0x10(%rcx)
140088e6a: 48 89 01                    	movq	%rax, (%rcx)
140088e6d: 4c 89 d0                    	movq	%r10, %rax
140088e70: 48 c7 41 08 00 00 00 00     	movq	$0x0, 0x8(%rcx)
140088e78: 48 83 c4 78                 	addq	$0x78, %rsp
140088e7c: 5b                          	popq	%rbx
140088e7d: 5e                          	popq	%rsi
140088e7e: c3                          	retq
140088e7f: 90                          	nop
140088e80: 48 89 8c 24 90 00 00 00     	movq	%rcx, 0x90(%rsp)
140088e88: 48 8b 35 b1 d1 0e 00        	movq	0xed1b1(%rip), %rsi     # 0x140176040
140088e8f: b9 e9 fd 00 00              	movl	$0xfde9, %ecx           # imm = 0xFDE9
140088e94: 48 8d 5c 24 60              	leaq	0x60(%rsp), %rbx
140088e99: 48 c7 44 24 38 00 00 00 00  	movq	$0x0, 0x38(%rsp)
140088ea2: 48 c7 44 24 30 00 00 00 00  	movq	$0x0, 0x30(%rsp)
140088eab: c7 44 24 28 00 00 00 00     	movl	$0x0, 0x28(%rsp)
140088eb3: 48 c7 44 24 20 00 00 00 00  	movq	$0x0, 0x20(%rsp)
140088ebc: 4c 8b 02                    	movq	(%rdx), %r8
140088ebf: 48 89 94 24 98 00 00 00     	movq	%rdx, 0x98(%rsp)
140088ec7: 31 d2                       	xorl	%edx, %edx
140088ec9: ff d6                       	callq	*%rsi
140088ecb: 45 31 c0                    	xorl	%r8d, %r8d
140088ece: 48 8d 4c 24 50              	leaq	0x50(%rsp), %rcx
140088ed3: 48 89 5c 24 50              	movq	%rbx, 0x50(%rsp)
140088ed8: 48 63 d0                    	movslq	%eax, %rdx
140088edb: 89 44 24 4c                 	movl	%eax, 0x4c(%rsp)
140088edf: e8 4c 6c 03 00              	callq	0x1400bfb30 <.text+0xbeb30>
140088ee4: 8b 44 24 4c                 	movl	0x4c(%rsp), %eax
140088ee8: 31 d2                       	xorl	%edx, %edx
140088eea: b9 e9 fd 00 00              	movl	$0xfde9, %ecx           # imm = 0xFDE9
140088eef: 4c 8b 9c 24 98 00 00 00     	movq	0x98(%rsp), %r11
140088ef7: 48 c7 44 24 38 00 00 00 00  	movq	$0x0, 0x38(%rsp)
140088f00: 89 44 24 28                 	movl	%eax, 0x28(%rsp)
140088f04: 48 8b 44 24 50              	movq	0x50(%rsp), %rax
140088f09: 48 c7 44 24 30 00 00 00 00  	movq	$0x0, 0x30(%rsp)
140088f12: 48 89 44 24 20              	movq	%rax, 0x20(%rsp)
140088f17: 45 8b 4b 08                 	movl	0x8(%r11), %r9d
140088f1b: 4d 8b 03                    	movq	(%r11), %r8
140088f1e: ff d6                       	callq	*%rsi
140088f20: 4c 8b 94 24 90 00 00 00     	movq	0x90(%rsp), %r10
140088f28: 48 8b 44 24 50              	movq	0x50(%rsp), %rax
140088f2d: 48 8b 4c 24 58              	movq	0x58(%rsp), %rcx
140088f32: 49 8d 52 10                 	leaq	0x10(%r10), %rdx
140088f36: 49 89 12                    	movq	%rdx, (%r10)
140088f39: 48 39 d8                    	cmpq	%rbx, %rax
140088f3c: 74 22                       	je	0x140088f60 <.text+0x87f60>
140088f3e: 49 89 02                    	movq	%rax, (%r10)
140088f41: 48 8b 44 24 60              	movq	0x60(%rsp), %rax
140088f46: 49 89 42 10                 	movq	%rax, 0x10(%r10)
140088f4a: 4c 89 d0                    	movq	%r10, %rax
140088f4d: 49 89 4a 08                 	movq	%rcx, 0x8(%r10)
140088f51: 48 83 c4 78                 	addq	$0x78, %rsp
140088f55: 5b                          	popq	%rbx
140088f56: 5e                          	popq	%rsi
140088f57: c3                          	retq
140088f58: 0f 1f 84 00 00 00 00 00     	nopl	(%rax,%rax)
140088f60: 4c 8d 41 01                 	leaq	0x1(%rcx), %r8
140088f64: 41 83 f8 40                 	cmpl	$0x40, %r8d
140088f68: 73 4e                       	jae	0x140088fb8 <.text+0x87fb8>
140088f6a: 41 f6 c0 20                 	testb	$0x20, %r8b
140088f6e: 0f 85 ce 00 00 00           	jne	0x140089042 <.text+0x88042>
140088f74: 41 f6 c0 10                 	testb	$0x10, %r8b
140088f78: 0f 85 fa 00 00 00           	jne	0x140089078 <.text+0x88078>
140088f7e: 41 f6 c0 08                 	testb	$0x8, %r8b
140088f82: 0f 85 0e 01 00 00           	jne	0x140089096 <.text+0x88096>
140088f88: 41 f6 c0 04                 	testb	$0x4, %r8b
140088f8c: 0f 85 1d 01 00 00           	jne	0x1400890af <.text+0x880af>
140088f92: 45 85 c0                    	testl	%r8d, %r8d
140088f95: 74 b3                       	je	0x140088f4a <.text+0x87f4a>
140088f97: 0f b6 03                    	movzbl	(%rbx), %eax
140088f9a: 41 88 42 10                 	movb	%al, 0x10(%r10)
140088f9e: 41 f6 c0 02                 	testb	$0x2, %r8b
140088fa2: 74 a6                       	je	0x140088f4a <.text+0x87f4a>
140088fa4: 45 89 c0                    	movl	%r8d, %r8d
140088fa7: 42 0f b7 44 03 fe           	movzwl	-0x2(%rbx,%r8), %eax
140088fad: 66 42 89 44 02 fe           	movw	%ax, -0x2(%rdx,%r8)
140088fb3: eb 95                       	jmp	0x140088f4a <.text+0x87f4a>
140088fb5: 0f 1f 00                    	nopl	(%rax)
140088fb8: 45 89 c3                    	movl	%r8d, %r11d
140088fbb: 41 83 e8 01                 	subl	$0x1, %r8d
140088fbf: 4c 01 db                    	addq	%r11, %rbx
140088fc2: f3 0f 6f 43 c0              	movdqu	-0x40(%rbx), %xmm0
140088fc7: 42 0f 11 44 1a c0           	movups	%xmm0, -0x40(%rdx,%r11)
140088fcd: f3 0f 6f 43 d0              	movdqu	-0x30(%rbx), %xmm0
140088fd2: 42 0f 11 44 1a d0           	movups	%xmm0, -0x30(%rdx,%r11)
140088fd8: f3 0f 6f 43 e0              	movdqu	-0x20(%rbx), %xmm0
140088fdd: 42 0f 11 44 1a e0           	movups	%xmm0, -0x20(%rdx,%r11)
140088fe3: f3 0f 6f 43 f0              	movdqu	-0x10(%rbx), %xmm0
140088fe8: 42 0f 11 44 1a f0           	movups	%xmm0, -0x10(%rdx,%r11)
140088fee: 41 83 f8 40                 	cmpl	$0x40, %r8d
140088ff2: 0f 82 52 ff ff ff           	jb	0x140088f4a <.text+0x87f4a>
140088ff8: 41 83 e0 c0                 	andl	$-0x40, %r8d
140088ffc: 45 31 db                    	xorl	%r11d, %r11d
140088fff: 45 89 d9                    	movl	%r11d, %r9d
140089002: 41 83 c3 40                 	addl	$0x40, %r11d
140089006: f3 42 0f 6f 1c 08           	movdqu	(%rax,%r9), %xmm3
14008900c: f3 42 0f 6f 54 08 10        	movdqu	0x10(%rax,%r9), %xmm2
140089013: f3 42 0f 6f 4c 08 20        	movdqu	0x20(%rax,%r9), %xmm1
14008901a: f3 42 0f 6f 44 08 30        	movdqu	0x30(%rax,%r9), %xmm0
140089021: 42 0f 11 1c 0a              	movups	%xmm3, (%rdx,%r9)
140089026: 42 0f 11 54 0a 10           	movups	%xmm2, 0x10(%rdx,%r9)
14008902c: 42 0f 11 4c 0a 20           	movups	%xmm1, 0x20(%rdx,%r9)
140089032: 42 0f 11 44 0a 30           	movups	%xmm0, 0x30(%rdx,%r9)
140089038: 45 39 c3                    	cmpl	%r8d, %r11d
14008903b: 72 c2                       	jb	0x140088fff <.text+0x87fff>
14008903d: e9 08 ff ff ff              	jmp	0x140088f4a <.text+0x87f4a>
140089042: f3 0f 6f 03                 	movdqu	(%rbx), %xmm0
140089046: 45 89 c0                    	movl	%r8d, %r8d
140089049: 4a 8d 44 02 20              	leaq	0x20(%rdx,%r8), %rax
14008904e: 41 0f 11 42 10              	movups	%xmm0, 0x10(%r10)
140089053: f3 0f 6f 43 10              	movdqu	0x10(%rbx), %xmm0
140089058: 0f 11 42 10                 	movups	%xmm0, 0x10(%rdx)
14008905c: 4a 8d 54 03 20              	leaq	0x20(%rbx,%r8), %rdx
140089061: f3 0f 6f 42 c0              	movdqu	-0x40(%rdx), %xmm0
140089066: 0f 11 40 c0                 	movups	%xmm0, -0x40(%rax)
14008906a: f3 0f 6f 42 d0              	movdqu	-0x30(%rdx), %xmm0
14008906f: 0f 11 40 d0                 	movups	%xmm0, -0x30(%rax)
140089073: e9 d2 fe ff ff              	jmp	0x140088f4a <.text+0x87f4a>
140089078: f3 0f 6f 03                 	movdqu	(%rbx), %xmm0
14008907c: 45 89 c0                    	movl	%r8d, %r8d
14008907f: 41 0f 11 42 10              	movups	%xmm0, 0x10(%r10)
140089084: f3 42 0f 6f 44 03 f0        	movdqu	-0x10(%rbx,%r8), %xmm0
14008908b: 42 0f 11 44 02 f0           	movups	%xmm0, -0x10(%rdx,%r8)
140089091: e9 b4 fe ff ff              	jmp	0x140088f4a <.text+0x87f4a>
140089096: 48 8b 03                    	movq	(%rbx), %rax
140089099: 45 89 c0                    	movl	%r8d, %r8d
14008909c: 49 89 42 10                 	movq	%rax, 0x10(%r10)
1400890a0: 4a 8b 44 03 f8              	movq	-0x8(%rbx,%r8), %rax
1400890a5: 4a 89 44 02 f8              	movq	%rax, -0x8(%rdx,%r8)
1400890aa: e9 9b fe ff ff              	jmp	0x140088f4a <.text+0x87f4a>
1400890af: 8b 03                       	movl	(%rbx), %eax
1400890b1: 45 89 c0                    	movl	%r8d, %r8d
1400890b4: 41 89 42 10                 	movl	%eax, 0x10(%r10)
1400890b8: 42 8b 44 03 fc              	movl	-0x4(%rbx,%r8), %eax
1400890bd: 42 89 44 02 fc              	movl	%eax, -0x4(%rdx,%r8)
1400890c2: e9 83 fe ff ff              	jmp	0x140088f4a <.text+0x87f4a>
1400890c7: 48 89 c3                    	movq	%rax, %rbx
1400890ca: e9 91 ef 04 00              	jmp	0x1400d8060 <.text+0xd7060>
1400890cf: 90                          	nop
1400890d0: 56                          	pushq	%rsi
1400890d1: 53                          	pushq	%rbx
1400890d2: 48 83 ec 68                 	subq	$0x68, %rsp
1400890d6: 4c 8b 4a 08                 	movq	0x8(%rdx), %r9
1400890da: 49 89 ca                    	movq	%rcx, %r10
1400890dd: 4d 85 c9                    	testq	%r9, %r9
1400890e0: 75 26                       	jne	0x140089108 <.text+0x88108>
1400890e2: 48 8d 41 10                 	leaq	0x10(%rcx), %rax
1400890e6: 48 c7 41 08 00 00 00 00     	movq	$0x0, 0x8(%rcx)
1400890ee: 48 89 01                    	movq	%rax, (%rcx)
1400890f1: 31 c0                       	xorl	%eax, %eax
1400890f3: 66 89 41 10                 	movw	%ax, 0x10(%rcx)
1400890f7: 4c 89 d0                    	movq	%r10, %rax
1400890fa: 48 83 c4 68                 	addq	$0x68, %rsp
1400890fe: 5b                          	popq	%rbx
1400890ff: 5e                          	popq	%rsi
140089100: c3                          	retq
140089101: 0f 1f 80 00 00 00 00        	nopl	(%rax)
140089108: 48 89 8c 24 80 00 00 00     	movq	%rcx, 0x80(%rsp)
140089110: 48 8b 35 e1 cd 0e 00        	movq	0xecde1(%rip), %rsi     # 0x140175ef8
140089117: b9 e9 fd 00 00              	movl	$0xfde9, %ecx           # imm = 0xFDE9
14008911c: 48 8d 5c 24 50              	leaq	0x50(%rsp), %rbx
140089121: c7 44 24 28 00 00 00 00     	movl	$0x0, 0x28(%rsp)
140089129: 48 c7 44 24 20 00 00 00 00  	movq	$0x0, 0x20(%rsp)
140089132: 4c 8b 02                    	movq	(%rdx), %r8
140089135: 48 89 94 24 88 00 00 00     	movq	%rdx, 0x88(%rsp)
14008913d: 31 d2                       	xorl	%edx, %edx
14008913f: ff d6                       	callq	*%rsi
140089141: 45 31 c0                    	xorl	%r8d, %r8d
140089144: 48 8d 4c 24 40              	leaq	0x40(%rsp), %rcx
140089149: 48 89 5c 24 40              	movq	%rbx, 0x40(%rsp)
14008914e: 48 63 d0                    	movslq	%eax, %rdx
140089151: 89 44 24 3c                 	movl	%eax, 0x3c(%rsp)
140089155: e8 46 ae 03 00              	callq	0x1400c3fa0 <.text+0xc2fa0>
14008915a: 8b 44 24 3c                 	movl	0x3c(%rsp), %eax
14008915e: 31 d2                       	xorl	%edx, %edx
140089160: b9 e9 fd 00 00              	movl	$0xfde9, %ecx           # imm = 0xFDE9
140089165: 4c 8b 9c 24 88 00 00 00     	movq	0x88(%rsp), %r11
14008916d: 89 44 24 28                 	movl	%eax, 0x28(%rsp)
140089171: 48 8b 44 24 40              	movq	0x40(%rsp), %rax
140089176: 48 89 44 24 20              	movq	%rax, 0x20(%rsp)
14008917b: 45 8b 4b 08                 	movl	0x8(%r11), %r9d
14008917f: 4d 8b 03                    	movq	(%r11), %r8
140089182: ff d6                       	callq	*%rsi
140089184: 4c 8b 94 24 80 00 00 00     	movq	0x80(%rsp), %r10
14008918c: 48 8b 44 24 40              	movq	0x40(%rsp), %rax
140089191: 48 8b 4c 24 48              	movq	0x48(%rsp), %rcx
140089196: 49 8d 52 10                 	leaq	0x10(%r10), %rdx
14008919a: 49 89 12                    	movq	%rdx, (%r10)
14008919d: 48 39 d8                    	cmpq	%rbx, %rax
1400891a0: 74 1e                       	je	0x1400891c0 <.text+0x881c0>
1400891a2: 49 89 02                    	movq	%rax, (%r10)
1400891a5: 48 8b 44 24 50              	movq	0x50(%rsp), %rax
1400891aa: 49 89 42 10                 	movq	%rax, 0x10(%r10)
1400891ae: 4c 89 d0                    	movq	%r10, %rax
1400891b1: 49 89 4a 08                 	movq	%rcx, 0x8(%r10)
1400891b5: 48 83 c4 68                 	addq	$0x68, %rsp
1400891b9: 5b                          	popq	%rbx
1400891ba: 5e                          	popq	%rsi
1400891bb: c3                          	retq
1400891bc: 0f 1f 40 00                 	nopl	(%rax)
1400891c0: 4c 8d 44 09 02              	leaq	0x2(%rcx,%rcx), %r8
1400891c5: 41 83 f8 40                 	cmpl	$0x40, %r8d
1400891c9: 73 55                       	jae	0x140089220 <.text+0x88220>
1400891cb: 41 f6 c0 20                 	testb	$0x20, %r8b
1400891cf: 0f 85 d5 00 00 00           	jne	0x1400892aa <.text+0x882aa>
1400891d5: 41 f6 c0 10                 	testb	$0x10, %r8b
1400891d9: 0f 85 01 01 00 00           	jne	0x1400892e0 <.text+0x882e0>
1400891df: 41 f6 c0 08                 	testb	$0x8, %r8b
1400891e3: 0f 85 15 01 00 00           	jne	0x1400892fe <.text+0x882fe>
1400891e9: 41 f6 c0 04                 	testb	$0x4, %r8b
1400891ed: 0f 85 24 01 00 00           	jne	0x140089317 <.text+0x88317>
1400891f3: 45 85 c0                    	testl	%r8d, %r8d
1400891f6: 74 b6                       	je	0x1400891ae <.text+0x881ae>
1400891f8: 0f b6 03                    	movzbl	(%rbx), %eax
1400891fb: 41 88 42 10                 	movb	%al, 0x10(%r10)
1400891ff: 41 f6 c0 02                 	testb	$0x2, %r8b
140089203: 74 a9                       	je	0x1400891ae <.text+0x881ae>
140089205: 45 89 c0                    	movl	%r8d, %r8d
140089208: 42 0f b7 44 03 fe           	movzwl	-0x2(%rbx,%r8), %eax
14008920e: 66 42 89 44 02 fe           	movw	%ax, -0x2(%rdx,%r8)
140089214: eb 98                       	jmp	0x1400891ae <.text+0x881ae>
140089216: 66 2e 0f 1f 84 00 00 00 00 00       	nopw	%cs:(%rax,%rax)
140089220: 45 89 c3                    	movl	%r8d, %r11d
140089223: 41 83 e8 01                 	subl	$0x1, %r8d
140089227: 4c 01 db                    	addq	%r11, %rbx
14008922a: f3 0f 6f 43 c0              	movdqu	-0x40(%rbx), %xmm0
14008922f: 42 0f 11 44 1a c0           	movups	%xmm0, -0x40(%rdx,%r11)
140089235: f3 0f 6f 43 d0              	movdqu	-0x30(%rbx), %xmm0
14008923a: 42 0f 11 44 1a d0           	movups	%xmm0, -0x30(%rdx,%r11)
140089240: f3 0f 6f 43 e0              	movdqu	-0x20(%rbx), %xmm0
140089245: 42 0f 11 44 1a e0           	movups	%xmm0, -0x20(%rdx,%r11)
14008924b: f3 0f 6f 43 f0              	movdqu	-0x10(%rbx), %xmm0
140089250: 42 0f 11 44 1a f0           	movups	%xmm0, -0x10(%rdx,%r11)
140089256: 41 83 f8 40                 	cmpl	$0x40, %r8d
14008925a: 0f 82 4e ff ff ff           	jb	0x1400891ae <.text+0x881ae>
140089260: 41 83 e0 c0                 	andl	$-0x40, %r8d
140089264: 45 31 db                    	xorl	%r11d, %r11d
140089267: 45 89 d9                    	movl	%r11d, %r9d
14008926a: 41 83 c3 40                 	addl	$0x40, %r11d
14008926e: f3 42 0f 6f 1c 08           	movdqu	(%rax,%r9), %xmm3
140089274: f3 42 0f 6f 54 08 10        	movdqu	0x10(%rax,%r9), %xmm2
14008927b: f3 42 0f 6f 4c 08 20        	movdqu	0x20(%rax,%r9), %xmm1
140089282: f3 42 0f 6f 44 08 30        	movdqu	0x30(%rax,%r9), %xmm0
140089289: 42 0f 11 1c 0a              	movups	%xmm3, (%rdx,%r9)
14008928e: 42 0f 11 54 0a 10           	movups	%xmm2, 0x10(%rdx,%r9)
140089294: 42 0f 11 4c 0a 20           	movups	%xmm1, 0x20(%rdx,%r9)
14008929a: 42 0f 11 44 0a 30           	movups	%xmm0, 0x30(%rdx,%r9)
1400892a0: 45 39 c3                    	cmpl	%r8d, %r11d
1400892a3: 72 c2                       	jb	0x140089267 <.text+0x88267>
1400892a5: e9 04 ff ff ff              	jmp	0x1400891ae <.text+0x881ae>
1400892aa: f3 0f 6f 03                 	movdqu	(%rbx), %xmm0
1400892ae: 45 89 c0                    	movl	%r8d, %r8d
1400892b1: 4a 8d 44 02 20              	leaq	0x20(%rdx,%r8), %rax
1400892b6: 41 0f 11 42 10              	movups	%xmm0, 0x10(%r10)
1400892bb: f3 0f 6f 43 10              	movdqu	0x10(%rbx), %xmm0
1400892c0: 0f 11 42 10                 	movups	%xmm0, 0x10(%rdx)
1400892c4: 4a 8d 54 03 20              	leaq	0x20(%rbx,%r8), %rdx
1400892c9: f3 0f 6f 42 c0              	movdqu	-0x40(%rdx), %xmm0
1400892ce: 0f 11 40 c0                 	movups	%xmm0, -0x40(%rax)
1400892d2: f3 0f 6f 42 d0              	movdqu	-0x30(%rdx), %xmm0
1400892d7: 0f 11 40 d0                 	movups	%xmm0, -0x30(%rax)
1400892db: e9 ce fe ff ff              	jmp	0x1400891ae <.text+0x881ae>
1400892e0: f3 0f 6f 03                 	movdqu	(%rbx), %xmm0
1400892e4: 45 89 c0                    	movl	%r8d, %r8d
1400892e7: 41 0f 11 42 10              	movups	%xmm0, 0x10(%r10)
1400892ec: f3 42 0f 6f 44 03 f0        	movdqu	-0x10(%rbx,%r8), %xmm0
1400892f3: 42 0f 11 44 02 f0           	movups	%xmm0, -0x10(%rdx,%r8)
1400892f9: e9 b0 fe ff ff              	jmp	0x1400891ae <.text+0x881ae>
1400892fe: 48 8b 03                    	movq	(%rbx), %rax
140089301: 45 89 c0                    	movl	%r8d, %r8d
140089304: 49 89 42 10                 	movq	%rax, 0x10(%r10)
140089308: 4a 8b 44 03 f8              	movq	-0x8(%rbx,%r8), %rax
14008930d: 4a 89 44 02 f8              	movq	%rax, -0x8(%rdx,%r8)
140089312: e9 97 fe ff ff              	jmp	0x1400891ae <.text+0x881ae>
140089317: 8b 03                       	movl	(%rbx), %eax
140089319: 45 89 c0                    	movl	%r8d, %r8d
14008931c: 41 89 42 10                 	movl	%eax, 0x10(%r10)
140089320: 42 8b 44 03 fc              	movl	-0x4(%rbx,%r8), %eax
140089325: 42 89 44 02 fc              	movl	%eax, -0x4(%rdx,%r8)
14008932a: e9 7f fe ff ff              	jmp	0x1400891ae <.text+0x881ae>
14008932f: 48 89 c3                    	movq	%rax, %rbx
140089332: e9 3c ed 04 00              	jmp	0x1400d8073 <.text+0xd7073>
140089337: 66 0f 1f 84 00 00 00 00 00  	nopw	(%rax,%rax)
140089340: 41 54                       	pushq	%r12
140089342: 55                          	pushq	%rbp
140089343: 57                          	pushq	%rdi
140089344: 56                          	pushq	%rsi
140089345: 53                          	pushq	%rbx
140089346: 48 83 c4 80                 	addq	$-0x80, %rsp
14008934a: 48 89 d6                    	movq	%rdx, %rsi
14008934d: 48 8d 15 50 18 07 00        	leaq	0x71850(%rip), %rdx     # 0x1400faba4
140089354: 4c 89 c7                    	movq	%r8, %rdi
140089357: 48 89 cb                    	movq	%rcx, %rbx
14008935a: 4c 8d 42 01                 	leaq	0x1(%rdx), %r8
14008935e: 48 8d 4c 24 40              	leaq	0x40(%rsp), %rcx
140089363: 48 8d 6c 24 50              	leaq	0x50(%rsp), %rbp
140089368: 48 89 6c 24 40              	movq	%rbp, 0x40(%rsp)
14008936d: e8 ae f7 ff ff              	callq	0x140088b20 <.text+0x87b20>
140089372: 48 89 f9                    	movq	%rdi, %rcx
140089375: e8 e6 c2 01 00              	callq	0x1400a5660 <.text+0xa4660>
14008937a: 48 ba fe ff ff ff ff ff ff 7f       	movabsq	$0x7ffffffffffffffe, %rdx # imm = 0x7FFFFFFFFFFFFFFE
140089384: 48 2b 54 24 48              	subq	0x48(%rsp), %rdx
140089389: 49 89 c0                    	movq	%rax, %r8
14008938c: 48 39 c2                    	cmpq	%rax, %rdx
14008938f: 0f 82 3a ed 04 00           	jb	0x1400d80cf <.text+0xd70cf>
140089395: 48 89 fa                    	movq	%rdi, %rdx
140089398: 48 8d 4c 24 40              	leaq	0x40(%rsp), %rcx
14008939d: e8 2e 95 03 00              	callq	0x1400c28d0 <.text+0xc18d0>
1400893a2: 4c 8d 64 24 70              	leaq	0x70(%rsp), %r12
1400893a7: 48 8d 48 10                 	leaq	0x10(%rax), %rcx
1400893ab: 4c 89 64 24 60              	movq	%r12, 0x60(%rsp)
1400893b0: 48 8b 10                    	movq	(%rax), %rdx
1400893b3: 48 39 ca                    	cmpq	%rcx, %rdx
1400893b6: 0f 84 e4 01 00 00           	je	0x1400895a0 <.text+0x885a0>
1400893bc: 48 89 54 24 60              	movq	%rdx, 0x60(%rsp)
1400893c1: 48 8b 50 10                 	movq	0x10(%rax), %rdx
1400893c5: 48 89 54 24 70              	movq	%rdx, 0x70(%rsp)
1400893ca: 48 8b 50 08                 	movq	0x8(%rax), %rdx
1400893ce: 48 89 08                    	movq	%rcx, (%rax)
1400893d1: c6 40 10 00                 	movb	$0x0, 0x10(%rax)
1400893d5: 48 c7 40 08 00 00 00 00     	movq	$0x0, 0x8(%rax)
1400893dd: 48 b8 03 00 00 00 00 00 00 80       	movabsq	$-0x7ffffffffffffffd, %rax # imm = 0x8000000000000003
1400893e7: 48 01 d0                    	addq	%rdx, %rax
1400893ea: 48 89 54 24 68              	movq	%rdx, 0x68(%rsp)
1400893ef: 48 83 f8 01                 	cmpq	$0x1, %rax
1400893f3: 0f 86 c0 ec 04 00           	jbe	0x1400d80b9 <.text+0xd70b9>
1400893f9: 48 8d 7c 24 60              	leaq	0x60(%rsp), %rdi
1400893fe: 41 b8 02 00 00 00           	movl	$0x2, %r8d
140089404: 48 8d 15 9b 17 07 00        	leaq	0x7179b(%rip), %rdx     # 0x1400faba6
14008940b: 48 89 f9                    	movq	%rdi, %rcx
14008940e: e8 bd 94 03 00              	callq	0x1400c28d0 <.text+0xc18d0>
140089413: 48 8d 7c 24 30              	leaq	0x30(%rsp), %rdi
140089418: 48 8d 50 10                 	leaq	0x10(%rax), %rdx
14008941c: 48 89 7c 24 20              	movq	%rdi, 0x20(%rsp)
140089421: 48 8b 08                    	movq	(%rax), %rcx
140089424: 48 39 d1                    	cmpq	%rdx, %rcx
140089427: 0f 84 43 01 00 00           	je	0x140089570 <.text+0x88570>
14008942d: 48 89 4c 24 20              	movq	%rcx, 0x20(%rsp)
140089432: 48 8b 48 10                 	movq	0x10(%rax), %rcx
140089436: 48 89 4c 24 30              	movq	%rcx, 0x30(%rsp)
14008943b: 48 8b 48 08                 	movq	0x8(%rax), %rcx
14008943f: 48 89 4c 24 28              	movq	%rcx, 0x28(%rsp)
140089444: 48 89 10                    	movq	%rdx, (%rax)
140089447: 48 8b 4c 24 60              	movq	0x60(%rsp), %rcx
14008944c: c6 40 10 00                 	movb	$0x0, 0x10(%rax)
140089450: 48 c7 40 08 00 00 00 00     	movq	$0x0, 0x8(%rax)
140089458: 4c 39 e1                    	cmpq	%r12, %rcx
14008945b: 74 0e                       	je	0x14008946b <.text+0x8846b>
14008945d: 48 8b 44 24 70              	movq	0x70(%rsp), %rax
140089462: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
140089466: e8 65 0c 04 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008946b: 48 8b 4c 24 40              	movq	0x40(%rsp), %rcx
140089470: 48 39 e9                    	cmpq	%rbp, %rcx
140089473: 74 0e                       	je	0x140089483 <.text+0x88483>
140089475: 48 8b 44 24 50              	movq	0x50(%rsp), %rax
14008947a: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008947e: e8 4d 0c 04 00              	callq	0x1400ca0d0 <.text+0xc90d0>
140089483: 4c 8b 4c 24 28              	movq	0x28(%rsp), %r9
140089488: 48 8b 54 24 20              	movq	0x20(%rsp), %rdx
14008948d: 45 31 c0                    	xorl	%r8d, %r8d
140089490: 48 89 f1                    	movq	%rsi, %rcx
140089493: e8 48 df 02 00              	callq	0x1400b73e0 <.text+0xb63e0>
140089498: 48 83 f8 ff                 	cmpq	$-0x1, %rax
14008949c: 74 29                       	je	0x1400894c7 <.text+0x884c7>
14008949e: 4c 8b 4e 08                 	movq	0x8(%rsi), %r9
1400894a2: 48 03 44 24 28              	addq	0x28(%rsp), %rax
1400894a7: 4c 39 c8                    	cmpq	%r9, %rax
1400894aa: 73 1b                       	jae	0x1400894c7 <.text+0x884c7>
1400894ac: 48 8b 0e                    	movq	(%rsi), %rcx
1400894af: 90                          	nop
1400894b0: 0f b6 14 01                 	movzbl	(%rcx,%rax), %edx
1400894b4: 80 fa 20                    	cmpb	$0x20, %dl
1400894b7: 74 05                       	je	0x1400894be <.text+0x884be>
1400894b9: 80 fa 09                    	cmpb	$0x9, %dl
1400894bc: 75 42                       	jne	0x140089500 <.text+0x88500>
1400894be: 48 83 c0 01                 	addq	$0x1, %rax
1400894c2: 4c 39 c8                    	cmpq	%r9, %rax
1400894c5: 75 e9                       	jne	0x1400894b0 <.text+0x884b0>
1400894c7: 48 8d 43 10                 	leaq	0x10(%rbx), %rax
1400894cb: c6 43 10 00                 	movb	$0x0, 0x10(%rbx)
1400894cf: 48 89 03                    	movq	%rax, (%rbx)
1400894d2: 48 c7 43 08 00 00 00 00     	movq	$0x0, 0x8(%rbx)
1400894da: 48 8b 4c 24 20              	movq	0x20(%rsp), %rcx
1400894df: 48 39 f9                    	cmpq	%rdi, %rcx
1400894e2: 74 0e                       	je	0x1400894f2 <.text+0x884f2>
1400894e4: 48 8b 44 24 30              	movq	0x30(%rsp), %rax
1400894e9: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
1400894ed: e8 de 0b 04 00              	callq	0x1400ca0d0 <.text+0xc90d0>
1400894f2: 48 89 d8                    	movq	%rbx, %rax
1400894f5: 48 83 ec 80                 	subq	$-0x80, %rsp
1400894f9: 5b                          	popq	%rbx
1400894fa: 5e                          	popq	%rsi
1400894fb: 5f                          	popq	%rdi
1400894fc: 5d                          	popq	%rbp
1400894fd: 41 5c                       	popq	%r12
1400894ff: c3                          	retq
140089500: 4c 39 c8                    	cmpq	%r9, %rax
140089503: 73 c2                       	jae	0x1400894c7 <.text+0x884c7>
140089505: 49 89 c0                    	movq	%rax, %r8
140089508: 80 3c 01 22                 	cmpb	$0x22, (%rcx,%rax)
14008950c: 75 1c                       	jne	0x14008952a <.text+0x8852a>
14008950e: e9 dd 00 00 00              	jmp	0x1400895f0 <.text+0x885f0>
140089513: 0f 1f 44 00 00              	nopl	(%rax,%rax)
140089518: 80 fa 09                    	cmpb	$0x9, %dl
14008951b: 0f 8f af 00 00 00           	jg	0x1400895d0 <.text+0x885d0>
140089521: 49 83 c0 01                 	addq	$0x1, %r8
140089525: 4d 39 c8                    	cmpq	%r9, %r8
140089528: 74 16                       	je	0x140089540 <.text+0x88540>
14008952a: 42 0f b6 14 01              	movzbl	(%rcx,%r8), %edx
14008952f: 80 fa 2c                    	cmpb	$0x2c, %dl
140089532: 7e e4                       	jle	0x140089518 <.text+0x88518>
140089534: 80 fa 7d                    	cmpb	$0x7d, %dl
140089537: 75 e8                       	jne	0x140089521 <.text+0x88521>
140089539: 0f 1f 80 00 00 00 00        	nopl	(%rax)
140089540: 48 8d 53 10                 	leaq	0x10(%rbx), %rdx
140089544: 49 29 c0                    	subq	%rax, %r8
140089547: 49 29 c1                    	subq	%rax, %r9
14008954a: 48 89 d9                    	movq	%rbx, %rcx
14008954d: 48 89 13                    	movq	%rdx, (%rbx)
140089550: 48 8b 16                    	movq	(%rsi), %rdx
140089553: 48 01 c2                    	addq	%rax, %rdx
140089556: 4d 39 c8                    	cmpq	%r9, %r8
140089559: 4d 0f 47 c1                 	cmovaq	%r9, %r8
14008955d: 49 01 d0                    	addq	%rdx, %r8
140089560: e8 bb f5 ff ff              	callq	0x140088b20 <.text+0x87b20>
140089565: e9 70 ff ff ff              	jmp	0x1400894da <.text+0x884da>
14008956a: 66 0f 1f 44 00 00           	nopw	(%rax,%rax)
140089570: 48 8b 48 08                 	movq	0x8(%rax), %rcx
140089574: 4c 8d 41 01                 	leaq	0x1(%rcx), %r8
140089578: 41 83 e0 3f                 	andl	$0x3f, %r8d
14008957c: 0f 84 bd fe ff ff           	je	0x14008943f <.text+0x8843f>
140089582: 31 c9                       	xorl	%ecx, %ecx
140089584: 41 89 c9                    	movl	%ecx, %r9d
140089587: 83 c1 01                    	addl	$0x1, %ecx
14008958a: 46 0f b6 14 0a              	movzbl	(%rdx,%r9), %r10d
14008958f: 46 88 14 0f                 	movb	%r10b, (%rdi,%r9)
140089593: 44 39 c1                    	cmpl	%r8d, %ecx
140089596: 72 ec                       	jb	0x140089584 <.text+0x88584>
140089598: e9 9e fe ff ff              	jmp	0x14008943b <.text+0x8843b>
14008959d: 0f 1f 00                    	nopl	(%rax)
1400895a0: 48 8b 50 08                 	movq	0x8(%rax), %rdx
1400895a4: 4c 8d 42 01                 	leaq	0x1(%rdx), %r8
1400895a8: 41 83 e0 3f                 	andl	$0x3f, %r8d
1400895ac: 0f 84 1c fe ff ff           	je	0x1400893ce <.text+0x883ce>
1400895b2: 31 d2                       	xorl	%edx, %edx
1400895b4: 41 89 d1                    	movl	%edx, %r9d
1400895b7: 83 c2 01                    	addl	$0x1, %edx
1400895ba: 46 0f b6 14 09              	movzbl	(%rcx,%r9), %r10d
1400895bf: 47 88 14 0c                 	movb	%r10b, (%r12,%r9)
1400895c3: 44 39 c2                    	cmpl	%r8d, %edx
1400895c6: 72 ec                       	jb	0x1400895b4 <.text+0x885b4>
1400895c8: e9 fd fd ff ff              	jmp	0x1400893ca <.text+0x883ca>
1400895cd: 0f 1f 00                    	nopl	(%rax)
1400895d0: 49 ba 00 24 00 00 01 10 00 00       	movabsq	$0x100100002400, %r10 # imm = 0x100100002400
1400895da: 49 0f a3 d2                 	btq	%rdx, %r10
1400895de: 0f 83 3d ff ff ff           	jae	0x140089521 <.text+0x88521>
1400895e4: e9 57 ff ff ff              	jmp	0x140089540 <.text+0x88540>
1400895e9: 0f 1f 80 00 00 00 00        	nopl	(%rax)
1400895f0: 4c 8d 40 01                 	leaq	0x1(%rax), %r8
1400895f4: 4c 89 c0                    	movq	%r8, %rax
1400895f7: 4d 39 c8                    	cmpq	%r9, %r8
1400895fa: 72 0c                       	jb	0x140089608 <.text+0x88608>
1400895fc: eb 5d                       	jmp	0x14008965b <.text+0x8865b>
1400895fe: 66 90                       	nop
140089600: 4c 89 d0                    	movq	%r10, %rax
140089603: 4c 39 c8                    	cmpq	%r9, %rax
140089606: 73 60                       	jae	0x140089668 <.text+0x88668>
140089608: 0f b6 14 01                 	movzbl	(%rcx,%rax), %edx
14008960c: 80 fa 22                    	cmpb	$0x22, %dl
14008960f: 74 45                       	je	0x140089656 <.text+0x88656>
140089611: 4c 8d 50 01                 	leaq	0x1(%rax), %r10
140089615: 80 fa 5c                    	cmpb	$0x5c, %dl
140089618: 75 e6                       	jne	0x140089600 <.text+0x88600>
14008961a: 4d 39 ca                    	cmpq	%r9, %r10
14008961d: 72 31                       	jb	0x140089650 <.text+0x88650>
14008961f: 4c 89 d0                    	movq	%r10, %rax
140089622: 4c 29 c0                    	subq	%r8, %rax
140089625: 48 8d 53 10                 	leaq	0x10(%rbx), %rdx
140089629: 4d 29 c1                    	subq	%r8, %r9
14008962c: 48 89 d9                    	movq	%rbx, %rcx
14008962f: 48 89 13                    	movq	%rdx, (%rbx)
140089632: 48 8b 16                    	movq	(%rsi), %rdx
140089635: 4c 01 c2                    	addq	%r8, %rdx
140089638: 49 39 c1                    	cmpq	%rax, %r9
14008963b: 4c 0f 47 c8                 	cmovaq	%rax, %r9
14008963f: 4e 8d 04 0a                 	leaq	(%rdx,%r9), %r8
140089643: e8 d8 f4 ff ff              	callq	0x140088b20 <.text+0x87b20>
140089648: e9 8d fe ff ff              	jmp	0x1400894da <.text+0x884da>
14008964d: 0f 1f 00                    	nopl	(%rax)
140089650: 48 83 c0 02                 	addq	$0x2, %rax
140089654: eb ad                       	jmp	0x140089603 <.text+0x88603>
140089656: 4c 29 c0                    	subq	%r8, %rax
140089659: eb ca                       	jmp	0x140089625 <.text+0x88625>
14008965b: 4d 39 c1                    	cmpq	%r8, %r9
14008965e: 0f 82 77 ea 04 00           	jb	0x1400d80db <.text+0xd70db>
140089664: 31 c0                       	xorl	%eax, %eax
140089666: eb bd                       	jmp	0x140089625 <.text+0x88625>
140089668: 49 89 c2                    	movq	%rax, %r10
14008966b: eb b2                       	jmp	0x14008961f <.text+0x8861f>
14008966d: e9 15 ea 04 00              	jmp	0x1400d8087 <.text+0xd7087>
140089672: e9 25 ea 04 00              	jmp	0x1400d809c <.text+0xd709c>
140089677: e9 4e ea 04 00              	jmp	0x1400d80ca <.text+0xd70ca>
14008967c: 0f 1f 40 00                 	nopl	(%rax)
140089680: 48 83 ec 58                 	subq	$0x58, %rsp
140089684: 4c 8d 05 6c 15 07 00        	leaq	0x7156c(%rip), %r8      # 0x1400fabf7
14008968b: 48 89 ca                    	movq	%rcx, %rdx
14008968e: 48 8d 4c 24 30              	leaq	0x30(%rsp), %rcx
140089693: e8 a8 fc ff ff              	callq	0x140089340 <.text+0x88340>
140089698: 48 8b 4c 24 30              	movq	0x30(%rsp), %rcx
14008969d: 48 8d 44 24 40              	leaq	0x40(%rsp), %rax
1400896a2: 48 83 7c 24 38 04           	cmpq	$0x4, 0x38(%rsp)
1400896a8: 74 26                       	je	0x1400896d0 <.text+0x886d0>
1400896aa: 48 39 c1                    	cmpq	%rax, %rcx
1400896ad: b8 00 00 00 00              	movl	$0x0, %eax
1400896b2: 74 17                       	je	0x1400896cb <.text+0x886cb>
1400896b4: 88 44 24 2f                 	movb	%al, 0x2f(%rsp)
1400896b8: 48 8b 44 24 40              	movq	0x40(%rsp), %rax
1400896bd: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
1400896c1: e8 0a 0a 04 00              	callq	0x1400ca0d0 <.text+0xc90d0>
1400896c6: 0f b6 44 24 2f              	movzbl	0x2f(%rsp), %eax
1400896cb: 48 83 c4 58                 	addq	$0x58, %rsp
1400896cf: c3                          	retq
1400896d0: 81 39 74 72 75 65           	cmpl	$0x65757274, (%rcx)     # imm = 0x65757274
1400896d6: 75 d2                       	jne	0x1400896aa <.text+0x886aa>
1400896d8: 48 39 c1                    	cmpq	%rax, %rcx
1400896db: b8 01 00 00 00              	movl	$0x1, %eax
1400896e0: 75 d2                       	jne	0x1400896b4 <.text+0x886b4>
1400896e2: eb e7                       	jmp	0x1400896cb <.text+0x886cb>
1400896e4: 90                          	nop
1400896e5: 66 66 2e 0f 1f 84 00 00 00 00 00    	nopw	%cs:(%rax,%rax)
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
