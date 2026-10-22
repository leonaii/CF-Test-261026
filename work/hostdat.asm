
CF上号器_v1.17.exe:	file format coff-x86-64

Disassembly of section .text:

0000000140001000 <.text>:
14008e000: 50                          	pushq	%rax
14008e001: 01 00                       	addl	%eax, (%rax)
14008e003: 00 4c 8b bc                 	addb	%cl, -0x44(%rbx,%rcx,4)
14008e007: 24 60                       	andb	$0x60, %al
14008e009: 01 00                       	addl	%eax, (%rax)
14008e00b: 00 48 89                    	addb	%cl, -0x77(%rax)
14008e00e: 44 24 58                    	andb	$0x58, %al
14008e011: 48 83 7c 24 58 00           	cmpq	$0x0, 0x58(%rsp)
14008e017: 0f 84 7a fb ff ff           	je	0x14008db97 <.text+0x8cb97>
14008e01d: e9 65 fb ff ff              	jmp	0x14008db87 <.text+0x8cb87>
14008e022: ba 20 00 00 00              	movl	$0x20, %edx
14008e027: 48 8d 8c 24 50 01 00 00     	leaq	0x150(%rsp), %rcx
14008e02f: 48 29 c2                    	subq	%rax, %rdx
14008e032: e8 c9 0f 03 00              	callq	0x1400bf000 <.text+0xbe000>
14008e037: e9 50 fd ff ff              	jmp	0x14008dd8c <.text+0x8cd8c>
14008e03c: 48 8d 15 2d d0 06 00        	leaq	0x6d02d(%rip), %rdx     # 0x1400fb070
14008e043: 48 89 f1                    	movq	%rsi, %rcx
14008e046: e8 35 f5 ff ff              	callq	0x14008d580 <.text+0x8c580>
14008e04b: 48 8b 84 24 50 01 00 00     	movq	0x150(%rsp), %rax
14008e053: 4c 8b bc 24 60 01 00 00     	movq	0x160(%rsp), %r15
14008e05b: 31 db                       	xorl	%ebx, %ebx
14008e05d: 48 89 44 24 58              	movq	%rax, 0x58(%rsp)
14008e062: eb ad                       	jmp	0x14008e011 <.text+0x8d011>
14008e064: 31 d2                       	xorl	%edx, %edx
14008e066: e8 7d 2c 00 00              	callq	0x140090ce8 <.text+0x8fce8>
14008e06b: 48 8d 15 50 d0 06 00        	leaq	0x6d050(%rip), %rdx     # 0x1400fb0c2
14008e072: 48 89 f1                    	movq	%rsi, %rcx
14008e075: e8 06 f5 ff ff              	callq	0x14008d580 <.text+0x8c580>
14008e07a: eb cf                       	jmp	0x14008e04b <.text+0x8d04b>
14008e07c: 48 8b 8c 24 a0 00 00 00     	movq	0xa0(%rsp), %rcx
14008e084: 31 d2                       	xorl	%edx, %edx
14008e086: e8 5d 2c 00 00              	callq	0x140090ce8 <.text+0x8fce8>
14008e08b: 48 8d 15 4a d0 06 00        	leaq	0x6d04a(%rip), %rdx     # 0x1400fb0dc
14008e092: 48 89 f1                    	movq	%rsi, %rcx
14008e095: e8 e6 f4 ff ff              	callq	0x14008d580 <.text+0x8c580>
14008e09a: 4c 8b bc 24 60 01 00 00     	movq	0x160(%rsp), %r15
14008e0a2: 31 db                       	xorl	%ebx, %ebx
14008e0a4: e9 68 ff ff ff              	jmp	0x14008e011 <.text+0x8d011>
14008e0a9: 48 8b 8c 24 a8 00 00 00     	movq	0xa8(%rsp), %rcx
14008e0b1: e8 0a 2c 00 00              	callq	0x140090cc0 <.text+0x8fcc0>
14008e0b6: 48 8b 8c 24 a0 00 00 00     	movq	0xa0(%rsp), %rcx
14008e0be: 31 d2                       	xorl	%edx, %edx
14008e0c0: e8 23 2c 00 00              	callq	0x140090ce8 <.text+0x8fce8>
14008e0c5: 48 8d 15 22 d0 06 00        	leaq	0x6d022(%rip), %rdx     # 0x1400fb0ee
14008e0cc: 48 89 f1                    	movq	%rsi, %rcx
14008e0cf: e8 ac f4 ff ff              	callq	0x14008d580 <.text+0x8c580>
14008e0d4: eb c4                       	jmp	0x14008e09a <.text+0x8d09a>
14008e0d6: 45 31 ff                    	xorl	%r15d, %r15d
14008e0d9: 45 31 e4                    	xorl	%r12d, %r12d
14008e0dc: 45 31 ed                    	xorl	%r13d, %r13d
14008e0df: e9 5d fe ff ff              	jmp	0x14008df41 <.text+0x8cf41>
14008e0e4: 48 8d 15 1d d0 06 00        	leaq	0x6d01d(%rip), %rdx     # 0x1400fb108
14008e0eb: 48 89 f1                    	movq	%rsi, %rcx
14008e0ee: e8 8d f4 ff ff              	callq	0x14008d580 <.text+0x8c580>
14008e0f3: 31 db                       	xorl	%ebx, %ebx
14008e0f5: e9 ef fe ff ff              	jmp	0x14008dfe9 <.text+0x8cfe9>
14008e0fa: 48 89 c3                    	movq	%rax, %rbx
14008e0fd: e9 1c a5 04 00              	jmp	0x1400d861e <.text+0xd761e>
14008e102: 48 89 c3                    	movq	%rax, %rbx
14008e105: e9 1f a5 04 00              	jmp	0x1400d8629 <.text+0xd7629>
14008e10a: 48 89 c3                    	movq	%rax, %rbx
14008e10d: e9 2c a5 04 00              	jmp	0x1400d863e <.text+0xd763e>
14008e112: 90                          	nop
14008e113: 90                          	nop
14008e114: 90                          	nop
14008e115: 90                          	nop
14008e116: 90                          	nop
14008e117: 90                          	nop
14008e118: 90                          	nop
14008e119: 90                          	nop
14008e11a: 90                          	nop
14008e11b: 90                          	nop
14008e11c: 90                          	nop
14008e11d: 90                          	nop
14008e11e: 90                          	nop
14008e11f: 90                          	nop
14008e120: 66 2e 0f 1f 84 00 00 00 00 00       	nopw	%cs:(%rax,%rax)
14008e12a: 66 2e 0f 1f 84 00 00 00 00 00       	nopw	%cs:(%rax,%rax)
14008e134: 66 2e 0f 1f 84 00 00 00 00 00       	nopw	%cs:(%rax,%rax)
14008e13e: 66 90                       	nop
14008e140: 53                          	pushq	%rbx
14008e141: 48 83 ec 50                 	subq	$0x50, %rsp
14008e145: 48 89 54 24 68              	movq	%rdx, 0x68(%rsp)
14008e14a: 48 89 cb                    	movq	%rcx, %rbx
14008e14d: 48 8d 15 8c d4 06 00        	leaq	0x6d48c(%rip), %rdx     # 0x1400fb5e0
14008e154: 48 8d 4c 24 40              	leaq	0x40(%rsp), %rcx
14008e159: 44 89 44 24 70              	movl	%r8d, 0x70(%rsp)
14008e15e: 45 31 c0                    	xorl	%r8d, %r8d
14008e161: 4c 89 4c 24 78              	movq	%r9, 0x78(%rsp)
14008e166: 41 b9 08 00 00 00           	movl	$0x8, %r9d
14008e16c: 48 c7 44 24 40 00 00 00 00  	movq	$0x0, 0x40(%rsp)
14008e175: 48 c7 44 24 48 00 00 00 00  	movq	$0x0, 0x48(%rsp)
14008e17e: e8 dd 2a 00 00              	callq	0x140090c60 <.text+0x8fc60>
14008e183: 85 c0                       	testl	%eax, %eax
14008e185: b8 00 00 00 00              	movl	$0x0, %eax
14008e18a: 78 57                       	js	0x14008e1e3 <.text+0x8d1e3>
14008e18c: 48 89 5c 24 20              	movq	%rbx, 0x20(%rsp)
14008e191: 48 8b 4c 24 40              	movq	0x40(%rsp), %rcx
14008e196: 45 31 c9                    	xorl	%r9d, %r9d
14008e199: 45 31 c0                    	xorl	%r8d, %r8d
14008e19c: c7 44 24 30 00 00 00 00     	movl	$0x0, 0x30(%rsp)
14008e1a4: 48 8d 54 24 48              	leaq	0x48(%rsp), %rdx
14008e1a9: c7 44 24 28 20 00 00 00     	movl	$0x20, 0x28(%rsp)
14008e1b1: e8 2a 2b 00 00              	callq	0x140090ce0 <.text+0x8fce0>
14008e1b6: 89 c3                       	movl	%eax, %ebx
14008e1b8: 85 c0                       	testl	%eax, %eax
14008e1ba: 79 34                       	jns	0x14008e1f0 <.text+0x8d1f0>
14008e1bc: 48 8b 4c 24 48              	movq	0x48(%rsp), %rcx
14008e1c1: 48 85 c9                    	testq	%rcx, %rcx
14008e1c4: 74 05                       	je	0x14008e1cb <.text+0x8d1cb>
14008e1c6: e8 fd 2a 00 00              	callq	0x140090cc8 <.text+0x8fcc8>
14008e1cb: 48 8b 4c 24 40              	movq	0x40(%rsp), %rcx
14008e1d0: 48 85 c9                    	testq	%rcx, %rcx
14008e1d3: 74 07                       	je	0x14008e1dc <.text+0x8d1dc>
14008e1d5: 31 d2                       	xorl	%edx, %edx
14008e1d7: e8 0c 2b 00 00              	callq	0x140090ce8 <.text+0x8fce8>
14008e1dc: 89 d8                       	movl	%ebx, %eax
14008e1de: f7 d0                       	notl	%eax
14008e1e0: c1 e8 1f                    	shrl	$0x1f, %eax
14008e1e3: 48 83 c4 50                 	addq	$0x50, %rsp
14008e1e7: 5b                          	popq	%rbx
14008e1e8: c3                          	retq
14008e1e9: 0f 1f 80 00 00 00 00        	nopl	(%rax)
14008e1f0: 44 8b 44 24 70              	movl	0x70(%rsp), %r8d
14008e1f5: 48 8b 54 24 68              	movq	0x68(%rsp), %rdx
14008e1fa: 45 31 c9                    	xorl	%r9d, %r9d
14008e1fd: 48 8b 4c 24 48              	movq	0x48(%rsp), %rcx
14008e202: e8 69 2a 00 00              	callq	0x140090c70 <.text+0x8fc70>
14008e207: 89 c3                       	movl	%eax, %ebx
14008e209: 85 c0                       	testl	%eax, %eax
14008e20b: 78 af                       	js	0x14008e1bc <.text+0x8d1bc>
14008e20d: 48 8b 54 24 78              	movq	0x78(%rsp), %rdx
14008e212: 48 8b 4c 24 48              	movq	0x48(%rsp), %rcx
14008e217: 45 31 c9                    	xorl	%r9d, %r9d
14008e21a: 41 b8 20 00 00 00           	movl	$0x20, %r8d
14008e220: e8 73 2a 00 00              	callq	0x140090c98 <.text+0x8fc98>
14008e225: 89 c3                       	movl	%eax, %ebx
14008e227: eb 93                       	jmp	0x14008e1bc <.text+0x8d1bc>
14008e229: 0f 1f 80 00 00 00 00        	nopl	(%rax)
14008e230: 48 83 ec 58                 	subq	$0x58, %rsp
14008e234: 45 31 c9                    	xorl	%r9d, %r9d
14008e237: 45 31 c0                    	xorl	%r8d, %r8d
14008e23a: 48 89 4c 24 60              	movq	%rcx, 0x60(%rsp)
14008e23f: 48 8d 4c 24 40              	leaq	0x40(%rsp), %rcx
14008e244: 48 89 54 24 68              	movq	%rdx, 0x68(%rsp)
14008e249: 48 8d 15 90 d3 06 00        	leaq	0x6d390(%rip), %rdx     # 0x1400fb5e0
14008e250: 48 c7 44 24 40 00 00 00 00  	movq	$0x0, 0x40(%rsp)
14008e259: 48 c7 44 24 48 00 00 00 00  	movq	$0x0, 0x48(%rsp)
14008e262: e8 f9 29 00 00              	callq	0x140090c60 <.text+0x8fc60>
14008e267: 85 c0                       	testl	%eax, %eax
14008e269: 78 53                       	js	0x14008e2be <.text+0x8d2be>
14008e26b: 48 8b 4c 24 40              	movq	0x40(%rsp), %rcx
14008e270: 45 31 c9                    	xorl	%r9d, %r9d
14008e273: 45 31 c0                    	xorl	%r8d, %r8d
14008e276: 48 8d 54 24 48              	leaq	0x48(%rsp), %rdx
14008e27b: c7 44 24 30 00 00 00 00     	movl	$0x0, 0x30(%rsp)
14008e283: c7 44 24 28 00 00 00 00     	movl	$0x0, 0x28(%rsp)
14008e28b: 48 c7 44 24 20 00 00 00 00  	movq	$0x0, 0x20(%rsp)
14008e294: e8 47 2a 00 00              	callq	0x140090ce0 <.text+0x8fce0>
14008e299: 85 c0                       	testl	%eax, %eax
14008e29b: 79 2b                       	jns	0x14008e2c8 <.text+0x8d2c8>
14008e29d: 48 8b 4c 24 48              	movq	0x48(%rsp), %rcx
14008e2a2: 48 85 c9                    	testq	%rcx, %rcx
14008e2a5: 74 05                       	je	0x14008e2ac <.text+0x8d2ac>
14008e2a7: e8 1c 2a 00 00              	callq	0x140090cc8 <.text+0x8fcc8>
14008e2ac: 48 8b 4c 24 40              	movq	0x40(%rsp), %rcx
14008e2b1: 48 85 c9                    	testq	%rcx, %rcx
14008e2b4: 74 08                       	je	0x14008e2be <.text+0x8d2be>
14008e2b6: 31 d2                       	xorl	%edx, %edx
14008e2b8: e8 2b 2a 00 00              	callq	0x140090ce8 <.text+0x8fce8>
14008e2bd: 90                          	nop
14008e2be: 48 83 c4 58                 	addq	$0x58, %rsp
14008e2c2: c3                          	retq
14008e2c3: 0f 1f 44 00 00              	nopl	(%rax,%rax)
14008e2c8: 48 8b 54 24 60              	movq	0x60(%rsp), %rdx
14008e2cd: 48 8b 4c 24 48              	movq	0x48(%rsp), %rcx
14008e2d2: 45 31 c9                    	xorl	%r9d, %r9d
14008e2d5: 41 b8 2b 00 00 00           	movl	$0x2b, %r8d
14008e2db: e8 90 29 00 00              	callq	0x140090c70 <.text+0x8fc70>
14008e2e0: 85 c0                       	testl	%eax, %eax
14008e2e2: 78 b9                       	js	0x14008e29d <.text+0x8d29d>
14008e2e4: 48 8b 54 24 68              	movq	0x68(%rsp), %rdx
14008e2e9: 48 8b 4c 24 48              	movq	0x48(%rsp), %rcx
14008e2ee: 45 31 c9                    	xorl	%r9d, %r9d
14008e2f1: 41 b8 20 00 00 00           	movl	$0x20, %r8d
14008e2f7: e8 9c 29 00 00              	callq	0x140090c98 <.text+0x8fc98>
14008e2fc: eb 9f                       	jmp	0x14008e29d <.text+0x8d29d>
14008e2fe: 66 90                       	nop
14008e300: 41 57                       	pushq	%r15
14008e302: 41 56                       	pushq	%r14
14008e304: 41 55                       	pushq	%r13
14008e306: 41 54                       	pushq	%r12
14008e308: 55                          	pushq	%rbp
14008e309: 57                          	pushq	%rdi
14008e30a: 56                          	pushq	%rsi
14008e30b: 53                          	pushq	%rbx
14008e30c: 48 83 ec 78                 	subq	$0x78, %rsp
14008e310: be 0d 00 00 00              	movl	$0xd, %esi
14008e315: 41 bb 03 00 00 00           	movl	$0x3, %r11d
14008e31b: bb 01 00 00 00              	movl	$0x1, %ebx
14008e320: 41 ba 02 9a c8 e4           	movl	$0xe4c89a02, %r10d      # imm = 0xE4C89A02
14008e326: 4c 8d 3d 53 d4 06 00        	leaq	0x6d453(%rip), %r15     # 0x1400fb780
14008e32d: 4c 8d 25 5c d4 06 00        	leaq	0x6d45c(%rip), %r12     # 0x1400fb790
14008e334: 4c 8d 35 35 d4 06 00        	leaq	0x6d435(%rip), %r14     # 0x1400fb770
14008e33b: c6 44 24 20 92              	movb	$-0x6e, 0x20(%rsp)
14008e340: 49 89 cd                    	movq	%rcx, %r13
14008e343: 48 89 d7                    	movq	%rdx, %rdi
14008e346: 48 8d 6c 24 21              	leaq	0x21(%rsp), %rbp
14008e34b: 4c 89 84 24 d0 00 00 00     	movq	%r8, 0xd0(%rsp)
14008e353: 41 b8 02 00 00 00           	movl	$0x2, %r8d
14008e359: 0f 1f 80 00 00 00 00        	nopl	(%rax)
14008e360: 44 89 c0                    	movl	%r8d, %eax
14008e363: 44 89 d2                    	movl	%r10d, %edx
14008e366: 89 d9                       	movl	%ebx, %ecx
14008e368: 45 89 d8                    	movl	%r11d, %r8d
14008e36b: 45 8b 14 87                 	movl	(%r15,%rax,4), %r10d
14008e36f: 48 89 d8                    	movq	%rbx, %rax
14008e372: 83 e1 03                    	andl	$0x3, %ecx
14008e375: 41 83 e0 03                 	andl	$0x3, %r8d
14008e379: 83 e0 03                    	andl	$0x3, %eax
14008e37c: c1 e1 03                    	shll	$0x3, %ecx
14008e37f: 41 83 c3 01                 	addl	$0x1, %r11d
14008e383: 48 83 c5 01                 	addq	$0x1, %rbp
14008e387: 41 8b 04 84                 	movl	(%r12,%rax,4), %eax
14008e38b: 45 89 d1                    	movl	%r10d, %r9d
14008e38e: c1 ea 03                    	shrl	$0x3, %edx
14008e391: d3 e8                       	shrl	%cl, %eax
14008e393: 42 8d 0c c5 00 00 00 00     	leal	(,%r8,8), %ecx
14008e39b: 41 d3 e9                    	shrl	%cl, %r9d
14008e39e: 89 d9                       	movl	%ebx, %ecx
14008e3a0: 44 31 c8                    	xorl	%r9d, %eax
14008e3a3: f7 d1                       	notl	%ecx
14008e3a5: 45 89 c1                    	movl	%r8d, %r9d
14008e3a8: 47 8b 0c 8e                 	movl	(%r14,%r9,4), %r9d
14008e3ac: 83 e1 03                    	andl	$0x3, %ecx
14008e3af: c1 e1 03                    	shll	$0x3, %ecx
14008e3b2: 41 d3 e9                    	shrl	%cl, %r9d
14008e3b5: 44 31 c8                    	xorl	%r9d, %eax
14008e3b8: 89 c1                       	movl	%eax, %ecx
14008e3ba: 89 d8                       	movl	%ebx, %eax
14008e3bc: 0f af c3                    	imull	%ebx, %eax
14008e3bf: 83 f0 5a                    	xorl	$0x5a, %eax
14008e3c2: 01 f0                       	addl	%esi, %eax
14008e3c4: 83 c6 0d                    	addl	$0xd, %esi
14008e3c7: 01 c8                       	addl	%ecx, %eax
14008e3c9: 8d 4b ff                    	leal	-0x1(%rbx), %ecx
14008e3cc: 48 83 c3 01                 	addq	$0x1, %rbx
14008e3d0: 83 e1 03                    	andl	$0x3, %ecx
14008e3d3: 41 8b 0c 8c                 	movl	(%r12,%rcx,4), %ecx
14008e3d7: c1 e9 0b                    	shrl	$0xb, %ecx
14008e3da: 31 c8                       	xorl	%ecx, %eax
14008e3dc: 31 d0                       	xorl	%edx, %eax
14008e3de: 88 45 ff                    	movb	%al, -0x1(%rbp)
14008e3e1: 48 83 fb 20                 	cmpq	$0x20, %rbx
14008e3e5: 0f 85 75 ff ff ff           	jne	0x14008e360 <.text+0x8d360>
14008e3eb: 66 0f 6f 44 24 20           	movdqa	0x20(%rsp), %xmm0
14008e3f1: 4c 89 ea                    	movq	%r13, %rdx
14008e3f4: 48 8d 4c 24 40              	leaq	0x40(%rsp), %rcx
14008e3f9: 48 b8 61 65 73 2d 68 6f 73 74       	movabsq	$0x74736f682d736561, %rax # imm = 0x74736F682D736561
14008e403: 48 89 44 24 60              	movq	%rax, 0x60(%rsp)
14008e408: 0f 29 44 24 40              	movaps	%xmm0, 0x40(%rsp)
14008e40d: 66 0f 6f 44 24 30           	movdqa	0x30(%rsp), %xmm0
14008e413: c7 44 24 67 74 2d 76 32     	movl	$0x32762d74, 0x67(%rsp) # imm = 0x32762D74
14008e41b: 0f 29 44 24 50              	movaps	%xmm0, 0x50(%rsp)
14008e420: e8 0b fe ff ff              	callq	0x14008e230 <.text+0x8d230>
14008e425: 48 89 fa                    	movq	%rdi, %rdx
14008e428: 48 8d 4c 24 40              	leaq	0x40(%rsp), %rcx
14008e42d: 48 b8 78 6f 72 2d 68 6f 73 74       	movabsq	$0x74736f682d726f78, %rax # imm = 0x74736F682D726F78
14008e437: 48 89 44 24 60              	movq	%rax, 0x60(%rsp)
14008e43c: 48 8d 7c 24 20              	leaq	0x20(%rsp), %rdi
14008e441: c7 44 24 67 74 2d 76 32     	movl	$0x32762d74, 0x67(%rsp) # imm = 0x32762D74
14008e449: e8 e2 fd ff ff              	callq	0x14008e230 <.text+0x8d230>
14008e44e: 48 8b 94 24 d0 00 00 00     	movq	0xd0(%rsp), %rdx
14008e456: 48 8d 4c 24 40              	leaq	0x40(%rsp), %rcx
14008e45b: 48 b8 6d 61 63 2d 68 6f 73 74       	movabsq	$0x74736f682d63616d, %rax # imm = 0x74736F682D63616D
14008e465: 48 89 44 24 60              	movq	%rax, 0x60(%rsp)
14008e46a: c7 44 24 67 74 2d 76 32     	movl	$0x32762d74, 0x67(%rsp) # imm = 0x32762D74
14008e472: e8 b9 fd ff ff              	callq	0x14008e230 <.text+0x8d230>
14008e477: 31 c0                       	xorl	%eax, %eax
14008e479: 48 89 d9                    	movq	%rbx, %rcx
14008e47c: f3 aa                       	rep		stosb	%al, %es:(%rdi)
14008e47e: b9 30 00 00 00              	movl	$0x30, %ecx
14008e483: 48 8d 7c 24 40              	leaq	0x40(%rsp), %rdi
14008e488: f3 aa                       	rep		stosb	%al, %es:(%rdi)
14008e48a: 48 83 c4 78                 	addq	$0x78, %rsp
14008e48e: 5b                          	popq	%rbx
14008e48f: 5e                          	popq	%rsi
14008e490: 5f                          	popq	%rdi
14008e491: 5d                          	popq	%rbp
14008e492: 41 5c                       	popq	%r12
14008e494: 41 5d                       	popq	%r13
14008e496: 41 5e                       	popq	%r14
14008e498: 41 5f                       	popq	%r15
14008e49a: c3                          	retq
14008e49b: 90                          	nop
14008e49c: 0f 1f 40 00                 	nopl	(%rax)
14008e4a0: 53                          	pushq	%rbx
14008e4a1: 48 83 ec 40                 	subq	$0x40, %rsp
14008e4a5: 49 29 d0                    	subq	%rdx, %r8
14008e4a8: 48 89 cb                    	movq	%rcx, %rbx
14008e4ab: 49 89 d2                    	movq	%rdx, %r10
14008e4ae: 4c 89 44 24 38              	movq	%r8, 0x38(%rsp)
14008e4b3: 4d 89 c1                    	movq	%r8, %r9
14008e4b6: 49 83 f8 0f                 	cmpq	$0xf, %r8
14008e4ba: 77 44                       	ja	0x14008e500 <.text+0x8d500>
14008e4bc: 48 8b 09                    	movq	(%rcx), %rcx
14008e4bf: 49 83 f8 01                 	cmpq	$0x1, %r8
14008e4c3: 75 23                       	jne	0x14008e4e8 <.text+0x8d4e8>
14008e4c5: 0f b6 02                    	movzbl	(%rdx), %eax
14008e4c8: 88 01                       	movb	%al, (%rcx)
14008e4ca: 4c 8b 4c 24 38              	movq	0x38(%rsp), %r9
14008e4cf: 48 8b 0b                    	movq	(%rbx), %rcx
14008e4d2: 42 c6 04 09 00              	movb	$0x0, (%rcx,%r9)
14008e4d7: 4c 89 4b 08                 	movq	%r9, 0x8(%rbx)
14008e4db: 48 83 c4 40                 	addq	$0x40, %rsp
14008e4df: 5b                          	popq	%rbx
14008e4e0: c3                          	retq
14008e4e1: 0f 1f 80 00 00 00 00        	nopl	(%rax)
14008e4e8: 4d 85 c0                    	testq	%r8, %r8
14008e4eb: 75 43                       	jne	0x14008e530 <.text+0x8d530>
14008e4ed: 42 c6 04 09 00              	movb	$0x0, (%rcx,%r9)
14008e4f2: 4c 89 4b 08                 	movq	%r9, 0x8(%rbx)
14008e4f6: 48 83 c4 40                 	addq	$0x40, %rsp
14008e4fa: 5b                          	popq	%rbx
14008e4fb: c3                          	retq
14008e4fc: 0f 1f 40 00                 	nopl	(%rax)
14008e500: 48 89 54 24 58              	movq	%rdx, 0x58(%rsp)
14008e505: 48 8d 54 24 38              	leaq	0x38(%rsp), %rdx
14008e50a: 4c 89 44 24 28              	movq	%r8, 0x28(%rsp)
14008e50f: 45 31 c0                    	xorl	%r8d, %r8d
14008e512: e8 79 45 03 00              	callq	0x1400c2a90 <.text+0xc1a90>
14008e517: 4c 8b 4c 24 28              	movq	0x28(%rsp), %r9
14008e51c: 4c 8b 54 24 58              	movq	0x58(%rsp), %r10
14008e521: 48 89 03                    	movq	%rax, (%rbx)
14008e524: 48 89 c1                    	movq	%rax, %rcx
14008e527: 48 8b 44 24 38              	movq	0x38(%rsp), %rax
14008e52c: 48 89 43 10                 	movq	%rax, 0x10(%rbx)
14008e530: 4d 89 c8                    	movq	%r9, %r8
14008e533: 4c 89 d2                    	movq	%r10, %rdx
14008e536: e8 ad 72 01 00              	callq	0x1400a57e8 <.text+0xa47e8>
14008e53b: 4c 8b 4c 24 38              	movq	0x38(%rsp), %r9
14008e540: 48 8b 0b                    	movq	(%rbx), %rcx
14008e543: 42 c6 04 09 00              	movb	$0x0, (%rcx,%r9)
14008e548: 4c 89 4b 08                 	movq	%r9, 0x8(%rbx)
14008e54c: 48 83 c4 40                 	addq	$0x40, %rsp
14008e550: 5b                          	popq	%rbx
14008e551: c3                          	retq
14008e552: 0f 1f 00                    	nopl	(%rax)
14008e555: 66 66 2e 0f 1f 84 00 00 00 00 00    	nopw	%cs:(%rax,%rax)
14008e560: 53                          	pushq	%rbx
14008e561: 48 83 ec 40                 	subq	$0x40, %rsp
14008e565: 49 29 d0                    	subq	%rdx, %r8
14008e568: 48 89 cb                    	movq	%rcx, %rbx
14008e56b: 49 89 d2                    	movq	%rdx, %r10
14008e56e: 4d 89 c3                    	movq	%r8, %r11
14008e571: 4d 89 c1                    	movq	%r8, %r9
14008e574: 49 d1 fb                    	sarq	%r11
14008e577: 4c 89 5c 24 38              	movq	%r11, 0x38(%rsp)
14008e57c: 4c 89 d8                    	movq	%r11, %rax
14008e57f: 49 83 f8 0e                 	cmpq	$0xe, %r8
14008e583: 77 1b                       	ja	0x14008e5a0 <.text+0x8d5a0>
14008e585: 48 8b 09                    	movq	(%rcx), %rcx
14008e588: 49 83 fb 01                 	cmpq	$0x1, %r11
14008e58c: 75 4c                       	jne	0x14008e5da <.text+0x8d5da>
14008e58e: 0f b7 12                    	movzwl	(%rdx), %edx
14008e591: 66 89 11                    	movw	%dx, (%rcx)
14008e594: ba 02 00 00 00              	movl	$0x2, %edx
14008e599: eb 50                       	jmp	0x14008e5eb <.text+0x8d5eb>
14008e59b: 0f 1f 44 00 00              	nopl	(%rax,%rax)
14008e5a0: 48 89 54 24 58              	movq	%rdx, 0x58(%rsp)
14008e5a5: 48 8d 54 24 38              	leaq	0x38(%rsp), %rdx
14008e5aa: 4c 89 44 24 20              	movq	%r8, 0x20(%rsp)
14008e5af: 45 31 c0                    	xorl	%r8d, %r8d
14008e5b2: 4c 89 5c 24 28              	movq	%r11, 0x28(%rsp)
14008e5b7: e8 24 82 03 00              	callq	0x1400c67e0 <.text+0xc57e0>
14008e5bc: 4c 8b 4c 24 20              	movq	0x20(%rsp), %r9
14008e5c1: 4c 8b 54 24 58              	movq	0x58(%rsp), %r10
14008e5c6: 48 89 03                    	movq	%rax, (%rbx)
14008e5c9: 48 89 c1                    	movq	%rax, %rcx
14008e5cc: 48 8b 44 24 38              	movq	0x38(%rsp), %rax
14008e5d1: 4c 8b 5c 24 28              	movq	0x28(%rsp), %r11
14008e5d6: 48 89 43 10                 	movq	%rax, 0x10(%rbx)
14008e5da: 4d 85 db                    	testq	%r11, %r11
14008e5dd: 75 21                       	jne	0x14008e600 <.text+0x8d600>
14008e5df: 48 8b 44 24 38              	movq	0x38(%rsp), %rax
14008e5e4: 48 8b 0b                    	movq	(%rbx), %rcx
14008e5e7: 48 8d 14 00                 	leaq	(%rax,%rax), %rdx
14008e5eb: 45 31 c0                    	xorl	%r8d, %r8d
14008e5ee: 66 44 89 04 11              	movw	%r8w, (%rcx,%rdx)
14008e5f3: 48 89 43 08                 	movq	%rax, 0x8(%rbx)
14008e5f7: 48 83 c4 40                 	addq	$0x40, %rsp
14008e5fb: 5b                          	popq	%rbx
14008e5fc: c3                          	retq
14008e5fd: 0f 1f 00                    	nopl	(%rax)
14008e600: 4d 89 c8                    	movq	%r9, %r8
14008e603: 4c 89 d2                    	movq	%r10, %rdx
14008e606: e8 dd 71 01 00              	callq	0x1400a57e8 <.text+0xa47e8>
14008e60b: eb d2                       	jmp	0x14008e5df <.text+0x8d5df>
14008e60d: 90                          	nop
14008e60e: 66 90                       	nop
14008e610: 56                          	pushq	%rsi
14008e611: 53                          	pushq	%rbx
14008e612: 48 83 ec 38                 	subq	$0x38, %rsp
14008e616: 4c 8b 41 08                 	movq	0x8(%rcx), %r8
14008e61a: 4c 8b 51 10                 	movq	0x10(%rcx), %r10
14008e61e: 4d 39 d0                    	cmpq	%r10, %r8
14008e621: 74 15                       	je	0x14008e638 <.text+0x8d638>
14008e623: 0f b6 02                    	movzbl	(%rdx), %eax
14008e626: 41 88 00                    	movb	%al, (%r8)
14008e629: 48 83 41 08 01              	addq	$0x1, 0x8(%rcx)
14008e62e: 48 83 c4 38                 	addq	$0x38, %rsp
14008e632: 5b                          	popq	%rbx
14008e633: 5e                          	popq	%rsi
14008e634: c3                          	retq
14008e635: 0f 1f 00                    	nopl	(%rax)
14008e638: 48 8b 01                    	movq	(%rcx), %rax
14008e63b: 48 89 c3                    	movq	%rax, %rbx
14008e63e: 49 29 c0                    	subq	%rax, %r8
14008e641: 48 b8 ff ff ff ff ff ff ff 7f       	movabsq	$0x7fffffffffffffff, %rax # imm = 0x7FFFFFFFFFFFFFFF
14008e64b: 49 39 c0                    	cmpq	%rax, %r8
14008e64e: 0f 84 5a a0 04 00           	je	0x1400d86ae <.text+0xd76ae>
14008e654: 4d 85 c0                    	testq	%r8, %r8
14008e657: 41 bb 01 00 00 00           	movl	$0x1, %r11d
14008e65d: 48 89 4c 24 50              	movq	%rcx, 0x50(%rsp)
14008e662: 48 b8 ff ff ff ff ff ff ff 7f       	movabsq	$0x7fffffffffffffff, %rax # imm = 0x7FFFFFFFFFFFFFFF
14008e66c: 4d 0f 45 d8                 	cmovneq	%r8, %r11
14008e670: 48 89 54 24 58              	movq	%rdx, 0x58(%rsp)
14008e675: 4c 89 54 24 28              	movq	%r10, 0x28(%rsp)
14008e67a: 4b 8d 0c 03                 	leaq	(%r11,%r8), %rcx
14008e67e: 4c 89 44 24 20              	movq	%r8, 0x20(%rsp)
14008e683: 48 39 c1                    	cmpq	%rax, %rcx
14008e686: 48 0f 46 c1                 	cmovbeq	%rcx, %rax
14008e68a: 48 89 c1                    	movq	%rax, %rcx
14008e68d: 48 89 c6                    	movq	%rax, %rsi
14008e690: e8 4b ba 03 00              	callq	0x1400ca0e0 <.text+0xc90e0>
14008e695: 48 8b 54 24 58              	movq	0x58(%rsp), %rdx
14008e69a: 4c 8b 44 24 20              	movq	0x20(%rsp), %r8
14008e69f: 49 89 c3                    	movq	%rax, %r11
14008e6a2: 4c 8b 54 24 28              	movq	0x28(%rsp), %r10
14008e6a7: 4c 8b 4c 24 50              	movq	0x50(%rsp), %r9
14008e6ac: 0f b6 02                    	movzbl	(%rdx), %eax
14008e6af: 4d 85 c0                    	testq	%r8, %r8
14008e6b2: 43 88 04 03                 	movb	%al, (%r11,%r8)
14008e6b6: 74 1d                       	je	0x14008e6d5 <.text+0x8d6d5>
14008e6b8: 4c 89 d9                    	movq	%r11, %rcx
14008e6bb: 48 89 da                    	movq	%rbx, %rdx
14008e6be: e8 25 71 01 00              	callq	0x1400a57e8 <.text+0xa47e8>
14008e6c3: 4c 8b 4c 24 50              	movq	0x50(%rsp), %r9
14008e6c8: 4c 8b 54 24 28              	movq	0x28(%rsp), %r10
14008e6cd: 4c 8b 44 24 20              	movq	0x20(%rsp), %r8
14008e6d2: 49 89 c3                    	movq	%rax, %r11
14008e6d5: 4b 8d 44 03 01              	leaq	0x1(%r11,%r8), %rax
14008e6da: 48 89 d9                    	movq	%rbx, %rcx
14008e6dd: 48 85 db                    	testq	%rbx, %rbx
14008e6e0: 74 29                       	je	0x14008e70b <.text+0x8d70b>
14008e6e2: 4c 89 d2                    	movq	%r10, %rdx
14008e6e5: 4c 89 5c 24 28              	movq	%r11, 0x28(%rsp)
14008e6ea: 48 29 da                    	subq	%rbx, %rdx
14008e6ed: 4c 89 4c 24 50              	movq	%r9, 0x50(%rsp)
14008e6f2: 48 89 44 24 20              	movq	%rax, 0x20(%rsp)
14008e6f7: e8 d4 b9 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008e6fc: 4c 8b 5c 24 28              	movq	0x28(%rsp), %r11
14008e701: 4c 8b 4c 24 50              	movq	0x50(%rsp), %r9
14008e706: 48 8b 44 24 20              	movq	0x20(%rsp), %rax
14008e70b: 49 89 41 08                 	movq	%rax, 0x8(%r9)
14008e70f: 49 8d 04 33                 	leaq	(%r11,%rsi), %rax
14008e713: 4d 89 19                    	movq	%r11, (%r9)
14008e716: 49 89 41 10                 	movq	%rax, 0x10(%r9)
14008e71a: 48 83 c4 38                 	addq	$0x38, %rsp
14008e71e: 5b                          	popq	%rbx
14008e71f: 5e                          	popq	%rsi
14008e720: c3                          	retq
14008e721: 90                          	nop
14008e722: 0f 1f 00                    	nopl	(%rax)
14008e725: 66 66 2e 0f 1f 84 00 00 00 00 00    	nopw	%cs:(%rax,%rax)
14008e730: 41 57                       	pushq	%r15
14008e732: 41 56                       	pushq	%r14
14008e734: 41 55                       	pushq	%r13
14008e736: 41 54                       	pushq	%r12
14008e738: 55                          	pushq	%rbp
14008e739: 57                          	pushq	%rdi
14008e73a: 56                          	pushq	%rsi
14008e73b: 53                          	pushq	%rbx
14008e73c: 48 83 ec 38                 	subq	$0x38, %rsp
14008e740: 48 89 cd                    	movq	%rcx, %rbp
14008e743: 48 89 d7                    	movq	%rdx, %rdi
14008e746: 4d 89 c4                    	movq	%r8, %r12
14008e749: 4d 39 c8                    	cmpq	%r9, %r8
14008e74c: 0f 84 a1 01 00 00           	je	0x14008e8f3 <.text+0x8d8f3>
14008e752: 4c 8b 71 10                 	movq	0x10(%rcx), %r14
14008e756: 48 8b 59 08                 	movq	0x8(%rcx), %rbx
14008e75a: 4c 89 ce                    	movq	%r9, %rsi
14008e75d: 4c 29 c6                    	subq	%r8, %rsi
14008e760: 4c 89 f0                    	movq	%r14, %rax
14008e763: 48 29 d8                    	subq	%rbx, %rax
14008e766: 48 39 f0                    	cmpq	%rsi, %rax
14008e769: 0f 82 91 00 00 00           	jb	0x14008e800 <.text+0x8d800>
14008e76f: 49 89 dd                    	movq	%rbx, %r13
14008e772: 49 29 d5                    	subq	%rdx, %r13
14008e775: 4c 39 ee                    	cmpq	%r13, %rsi
14008e778: 0f 83 8a 01 00 00           	jae	0x14008e908 <.text+0x8d908>
14008e77e: 4c 89 c2                    	movq	%r8, %rdx
14008e781: 45 31 ed                    	xorl	%r13d, %r13d
14008e784: 4c 29 ca                    	subq	%r9, %rdx
14008e787: 48 01 da                    	addq	%rbx, %rdx
14008e78a: 66 66 2e 0f 1f 84 00 00 00 00 00    	nopw	%cs:(%rax,%rax)
14008e795: 66 66 2e 0f 1f 84 00 00 00 00 00    	nopw	%cs:(%rax,%rax)
14008e7a0: 42 0f b6 04 2a              	movzbl	(%rdx,%r13), %eax
14008e7a5: 42 88 04 2b                 	movb	%al, (%rbx,%r13)
14008e7a9: 49 83 c5 01                 	addq	$0x1, %r13
14008e7ad: 4c 39 ee                    	cmpq	%r13, %rsi
14008e7b0: 75 ee                       	jne	0x14008e7a0 <.text+0x8d7a0>
14008e7b2: 49 89 d0                    	movq	%rdx, %r8
14008e7b5: 48 89 f9                    	movq	%rdi, %rcx
14008e7b8: 48 01 75 08                 	addq	%rsi, 0x8(%rbp)
14008e7bc: 48 29 d1                    	subq	%rdx, %rcx
14008e7bf: 49 29 f8                    	subq	%rdi, %r8
14008e7c2: 48 01 d9                    	addq	%rbx, %rcx
14008e7c5: 49 83 f8 01                 	cmpq	$0x1, %r8
14008e7c9: 0f 8e a1 01 00 00           	jle	0x14008e970 <.text+0x8d970>
14008e7cf: 48 89 fa                    	movq	%rdi, %rdx
14008e7d2: e8 19 70 01 00              	callq	0x1400a57f0 <.text+0xa47f0>
14008e7d7: 48 83 fe 01                 	cmpq	$0x1, %rsi
14008e7db: 0f 8e 84 01 00 00           	jle	0x14008e965 <.text+0x8d965>
14008e7e1: 4d 89 e8                    	movq	%r13, %r8
14008e7e4: 4c 89 e2                    	movq	%r12, %rdx
14008e7e7: 48 89 f9                    	movq	%rdi, %rcx
14008e7ea: 48 83 c4 38                 	addq	$0x38, %rsp
14008e7ee: 5b                          	popq	%rbx
14008e7ef: 5e                          	popq	%rsi
14008e7f0: 5f                          	popq	%rdi
14008e7f1: 5d                          	popq	%rbp
14008e7f2: 41 5c                       	popq	%r12
14008e7f4: 41 5d                       	popq	%r13
14008e7f6: 41 5e                       	popq	%r14
14008e7f8: 41 5f                       	popq	%r15
14008e7fa: e9 f1 6f 01 00              	jmp	0x1400a57f0 <.text+0xa47f0>
14008e7ff: 90                          	nop
14008e800: 48 b8 ff ff ff ff ff ff ff 7f       	movabsq	$0x7fffffffffffffff, %rax # imm = 0x7FFFFFFFFFFFFFFF
14008e80a: 4c 8b 29                    	movq	(%rcx), %r13
14008e80d: 49 89 da                    	movq	%rbx, %r10
14008e810: 4d 29 ea                    	subq	%r13, %r10
14008e813: 4c 29 d0                    	subq	%r10, %rax
14008e816: 48 39 f0                    	cmpq	%rsi, %rax
14008e819: 0f 82 9d 9e 04 00           	jb	0x1400d86bc <.text+0xd76bc>
14008e81f: 4c 39 d6                    	cmpq	%r10, %rsi
14008e822: 4c 89 d0                    	movq	%r10, %rax
14008e825: 4c 89 54 24 20              	movq	%r10, 0x20(%rsp)
14008e82a: 48 0f 43 c6                 	cmovaeq	%rsi, %rax
14008e82e: 4a 8d 0c 10                 	leaq	(%rax,%r10), %rcx
14008e832: 48 b8 ff ff ff ff ff ff ff 7f       	movabsq	$0x7fffffffffffffff, %rax # imm = 0x7FFFFFFFFFFFFFFF
14008e83c: 48 39 c1                    	cmpq	%rax, %rcx
14008e83f: 48 0f 46 c1                 	cmovbeq	%rcx, %rax
14008e843: 48 89 c1                    	movq	%rax, %rcx
14008e846: 49 89 c7                    	movq	%rax, %r15
14008e849: e8 92 b8 03 00              	callq	0x1400ca0e0 <.text+0xc90e0>
14008e84e: 49 89 c1                    	movq	%rax, %r9
14008e851: 4c 39 ef                    	cmpq	%r13, %rdi
14008e854: 0f 84 26 01 00 00           	je	0x14008e980 <.text+0x8d980>
14008e85a: 48 89 f9                    	movq	%rdi, %rcx
14008e85d: 31 c0                       	xorl	%eax, %eax
14008e85f: 4c 29 e9                    	subq	%r13, %rcx
14008e862: 0f 1f 84 00 00 00 00 00     	nopl	(%rax,%rax)
14008e86a: 66 66 2e 0f 1f 84 00 00 00 00 00    	nopw	%cs:(%rax,%rax)
14008e875: 66 66 2e 0f 1f 84 00 00 00 00 00    	nopw	%cs:(%rax,%rax)
14008e880: 41 0f b6 54 05 00           	movzbl	(%r13,%rax), %edx
14008e886: 41 88 14 01                 	movb	%dl, (%r9,%rax)
14008e88a: 48 83 c0 01                 	addq	$0x1, %rax
14008e88e: 48 39 c1                    	cmpq	%rax, %rcx
14008e891: 75 ed                       	jne	0x14008e880 <.text+0x8d880>
14008e893: 4c 01 c9                    	addq	%r9, %rcx
14008e896: 49 89 f0                    	movq	%rsi, %r8
14008e899: 4c 89 e2                    	movq	%r12, %rdx
14008e89c: 4c 89 4c 24 20              	movq	%r9, 0x20(%rsp)
14008e8a1: e8 42 6f 01 00              	callq	0x1400a57e8 <.text+0xa47e8>
14008e8a6: 4c 8b 4c 24 20              	movq	0x20(%rsp), %r9
14008e8ab: 48 01 c6                    	addq	%rax, %rsi
14008e8ae: 48 39 df                    	cmpq	%rbx, %rdi
14008e8b1: 74 19                       	je	0x14008e8cc <.text+0x8d8cc>
14008e8b3: 48 29 fb                    	subq	%rdi, %rbx
14008e8b6: 48 89 f1                    	movq	%rsi, %rcx
14008e8b9: 48 89 fa                    	movq	%rdi, %rdx
14008e8bc: 49 89 d8                    	movq	%rbx, %r8
14008e8bf: 48 01 de                    	addq	%rbx, %rsi
14008e8c2: e8 21 6f 01 00              	callq	0x1400a57e8 <.text+0xa47e8>
14008e8c7: 4c 8b 4c 24 20              	movq	0x20(%rsp), %r9
14008e8cc: 4c 89 f2                    	movq	%r14, %rdx
14008e8cf: 4c 89 e9                    	movq	%r13, %rcx
14008e8d2: 4c 89 4c 24 20              	movq	%r9, 0x20(%rsp)
14008e8d7: 4c 29 ea                    	subq	%r13, %rdx
14008e8da: e8 f1 b7 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008e8df: 4c 8b 4c 24 20              	movq	0x20(%rsp), %r9
14008e8e4: 4c 89 4d 00                 	movq	%r9, (%rbp)
14008e8e8: 4d 01 f9                    	addq	%r15, %r9
14008e8eb: 48 89 75 08                 	movq	%rsi, 0x8(%rbp)
14008e8ef: 4c 89 4d 10                 	movq	%r9, 0x10(%rbp)
14008e8f3: 48 83 c4 38                 	addq	$0x38, %rsp
14008e8f7: 5b                          	popq	%rbx
14008e8f8: 5e                          	popq	%rsi
14008e8f9: 5f                          	popq	%rdi
14008e8fa: 5d                          	popq	%rbp
14008e8fb: 41 5c                       	popq	%r12
14008e8fd: 41 5d                       	popq	%r13
14008e8ff: 41 5e                       	popq	%r14
14008e901: 41 5f                       	popq	%r15
14008e903: c3                          	retq
14008e904: 0f 1f 40 00                 	nopl	(%rax)
14008e908: 4b 8d 14 28                 	leaq	(%r8,%r13), %rdx
14008e90c: 4d 89 c8                    	movq	%r9, %r8
14008e90f: 49 29 d0                    	subq	%rdx, %r8
14008e912: 4d 85 c0                    	testq	%r8, %r8
14008e915: 0f 8e ab 00 00 00           	jle	0x14008e9c6 <.text+0x8d9c6>
14008e91b: 48 89 d9                    	movq	%rbx, %rcx
14008e91e: e8 c5 6e 01 00              	callq	0x1400a57e8 <.text+0xa47e8>
14008e923: 48 8b 45 08                 	movq	0x8(%rbp), %rax
14008e927: 4c 29 ee                    	subq	%r13, %rsi
14008e92a: 48 01 c6                    	addq	%rax, %rsi
14008e92d: 48 89 75 08                 	movq	%rsi, 0x8(%rbp)
14008e931: 48 39 df                    	cmpq	%rbx, %rdi
14008e934: 74 1e                       	je	0x14008e954 <.text+0x8d954>
14008e936: 48 29 fb                    	subq	%rdi, %rbx
14008e939: 31 c0                       	xorl	%eax, %eax
14008e93b: 0f 1f 44 00 00              	nopl	(%rax,%rax)
14008e940: 0f b6 14 07                 	movzbl	(%rdi,%rax), %edx
14008e944: 88 14 06                    	movb	%dl, (%rsi,%rax)
14008e947: 48 83 c0 01                 	addq	$0x1, %rax
14008e94b: 48 39 d8                    	cmpq	%rbx, %rax
14008e94e: 75 f0                       	jne	0x14008e940 <.text+0x8d940>
14008e950: 48 8b 75 08                 	movq	0x8(%rbp), %rsi
14008e954: 4c 01 ee                    	addq	%r13, %rsi
14008e957: 48 89 75 08                 	movq	%rsi, 0x8(%rbp)
14008e95b: 49 83 fd 01                 	cmpq	$0x1, %r13
14008e95f: 0f 8f 7c fe ff ff           	jg	0x14008e7e1 <.text+0x8d7e1>
14008e965: 75 8c                       	jne	0x14008e8f3 <.text+0x8d8f3>
14008e967: 41 0f b6 04 24              	movzbl	(%r12), %eax
14008e96c: 88 07                       	movb	%al, (%rdi)
14008e96e: eb 83                       	jmp	0x14008e8f3 <.text+0x8d8f3>
14008e970: 0f 85 61 fe ff ff           	jne	0x14008e7d7 <.text+0x8d7d7>
14008e976: 0f b6 07                    	movzbl	(%rdi), %eax
14008e979: 88 01                       	movb	%al, (%rcx)
14008e97b: e9 57 fe ff ff              	jmp	0x14008e7d7 <.text+0x8d7d7>
14008e980: 49 89 f0                    	movq	%rsi, %r8
14008e983: 4c 89 e2                    	movq	%r12, %rdx
14008e986: 48 89 c1                    	movq	%rax, %rcx
14008e989: e8 5a 6e 01 00              	callq	0x1400a57e8 <.text+0xa47e8>
14008e98e: 49 89 c1                    	movq	%rax, %r9
14008e991: 48 01 c6                    	addq	%rax, %rsi
14008e994: 49 39 dd                    	cmpq	%rbx, %r13
14008e997: 74 1f                       	je	0x14008e9b8 <.text+0x8d9b8>
14008e999: 4c 8b 44 24 20              	movq	0x20(%rsp), %r8
14008e99e: 48 89 f1                    	movq	%rsi, %rcx
14008e9a1: 4c 89 ea                    	movq	%r13, %rdx
14008e9a4: 48 89 44 24 28              	movq	%rax, 0x28(%rsp)
14008e9a9: e8 3a 6e 01 00              	callq	0x1400a57e8 <.text+0xa47e8>
14008e9ae: 4c 8b 4c 24 28              	movq	0x28(%rsp), %r9
14008e9b3: 48 03 74 24 20              	addq	0x20(%rsp), %rsi
14008e9b8: 4d 85 ed                    	testq	%r13, %r13
14008e9bb: 0f 84 23 ff ff ff           	je	0x14008e8e4 <.text+0x8d8e4>
14008e9c1: e9 06 ff ff ff              	jmp	0x14008e8cc <.text+0x8d8cc>
14008e9c6: 48 89 d8                    	movq	%rbx, %rax
14008e9c9: e9 59 ff ff ff              	jmp	0x14008e927 <.text+0x8d927>
14008e9ce: 66 90                       	nop
14008e9d0: 41 56                       	pushq	%r14
14008e9d2: 57                          	pushq	%rdi
14008e9d3: 56                          	pushq	%rsi
14008e9d4: 53                          	pushq	%rbx
14008e9d5: 48 83 ec 78                 	subq	$0x78, %rsp
14008e9d9: 45 31 c0                    	xorl	%r8d, %r8d
14008e9dc: 48 8b 41 08                 	movq	0x8(%rcx), %rax
14008e9e0: 48 83 e8 0a                 	subq	$0xa, %rax
14008e9e4: 48 83 f8 56                 	cmpq	$0x56, %rax
14008e9e8: 76 16                       	jbe	0x14008ea00 <.text+0x8da00>
14008e9ea: 44 89 c0                    	movl	%r8d, %eax
14008e9ed: 48 83 c4 78                 	addq	$0x78, %rsp
14008e9f1: 5b                          	popq	%rbx
14008e9f2: 5e                          	popq	%rsi
14008e9f3: 5f                          	popq	%rdi
14008e9f4: 41 5e                       	popq	%r14
14008e9f6: c3                          	retq
14008e9f7: 66 0f 1f 84 00 00 00 00 00  	nopw	(%rax,%rax)
14008ea00: 45 31 c0                    	xorl	%r8d, %r8d
14008ea03: 48 8d 5c 24 40              	leaq	0x40(%rsp), %rbx
14008ea08: 41 b9 07 00 00 00           	movl	$0x7, %r9d
14008ea0e: 48 89 8c 24 a0 00 00 00     	movq	%rcx, 0xa0(%rsp)
14008ea16: 48 8d 15 03 cc 06 00        	leaq	0x6cc03(%rip), %rdx     # 0x1400fb620
14008ea1d: 48 89 5c 24 30              	movq	%rbx, 0x30(%rsp)
14008ea22: c6 44 24 40 00              	movb	$0x0, 0x40(%rsp)
14008ea27: 48 c7 44 24 38 00 00 00 00  	movq	$0x0, 0x38(%rsp)
14008ea30: e8 8b 8b 02 00              	callq	0x1400b75c0 <.text+0xb65c0>
14008ea35: 48 8b 8c 24 a0 00 00 00     	movq	0xa0(%rsp), %rcx
14008ea3d: 48 85 c0                    	testq	%rax, %rax
14008ea40: 0f 84 3a 01 00 00           	je	0x14008eb80 <.text+0x8db80>
14008ea46: 45 31 c0                    	xorl	%r8d, %r8d
14008ea49: 41 b9 08 00 00 00           	movl	$0x8, %r9d
14008ea4f: 48 8d 15 21 cc 06 00        	leaq	0x6cc21(%rip), %rdx     # 0x1400fb677
14008ea56: 48 89 8c 24 a0 00 00 00     	movq	%rcx, 0xa0(%rsp)
14008ea5e: e8 5d 8b 02 00              	callq	0x1400b75c0 <.text+0xb65c0>
14008ea63: 48 85 c0                    	testq	%rax, %rax
14008ea66: 0f 85 da 00 00 00           	jne	0x14008eb46 <.text+0x8db46>
14008ea6c: 48 8b 8c 24 a0 00 00 00     	movq	0xa0(%rsp), %rcx
14008ea74: 4c 8b 49 08                 	movq	0x8(%rcx), %r9
14008ea78: 49 83 f9 07                 	cmpq	$0x7, %r9
14008ea7c: 0f 86 8b 9c 04 00           	jbe	0x1400d870d <.text+0xd770d>
14008ea82: 48 8d 44 24 60              	leaq	0x60(%rsp), %rax
14008ea87: 4c 8d 74 24 50              	leaq	0x50(%rsp), %r14
14008ea8c: 48 89 44 24 50              	movq	%rax, 0x50(%rsp)
14008ea91: 48 89 c6                    	movq	%rax, %rsi
14008ea94: 48 8b 01                    	movq	(%rcx), %rax
14008ea97: 48 8d 4c 24 50              	leaq	0x50(%rsp), %rcx
14008ea9c: 48 8d 7c 24 30              	leaq	0x30(%rsp), %rdi
14008eaa1: 48 8d 50 08                 	leaq	0x8(%rax), %rdx
14008eaa5: 4e 8d 04 08                 	leaq	(%rax,%r9), %r8
14008eaa9: e8 f2 f9 ff ff              	callq	0x14008e4a0 <.text+0x8d4a0>
14008eaae: 48 8b 44 24 30              	movq	0x30(%rsp), %rax
14008eab3: 49 89 c0                    	movq	%rax, %r8
14008eab6: 48 39 d8                    	cmpq	%rbx, %rax
14008eab9: 0f 84 d9 02 00 00           	je	0x14008ed98 <.text+0x8dd98>
14008eabf: 48 8b 4c 24 50              	movq	0x50(%rsp), %rcx
14008eac4: 48 8b 54 24 58              	movq	0x58(%rsp), %rdx
14008eac9: 48 39 f1                    	cmpq	%rsi, %rcx
14008eacc: 0f 84 de 02 00 00           	je	0x14008edb0 <.text+0x8ddb0>
14008ead2: 4c 8b 44 24 40              	movq	0x40(%rsp), %r8
14008ead7: 66 48 0f 6e c2              	movq	%rdx, %xmm0
14008eadc: 48 89 4c 24 30              	movq	%rcx, 0x30(%rsp)
14008eae1: 0f 16 44 24 60              	movhps	0x60(%rsp), %xmm0       # xmm0 = xmm0[0,1],mem[0,1]
14008eae6: 0f 11 44 24 38              	movups	%xmm0, 0x38(%rsp)
14008eaeb: 48 85 c0                    	testq	%rax, %rax
14008eaee: 0f 84 a4 01 00 00           	je	0x14008ec98 <.text+0x8dc98>
14008eaf4: 48 89 44 24 50              	movq	%rax, 0x50(%rsp)
14008eaf9: 4c 89 44 24 60              	movq	%r8, 0x60(%rsp)
14008eafe: c6 00 00                    	movb	$0x0, (%rax)
14008eb01: 48 8b 4c 24 50              	movq	0x50(%rsp), %rcx
14008eb06: 48 39 f1                    	cmpq	%rsi, %rcx
14008eb09: 74 0e                       	je	0x14008eb19 <.text+0x8db19>
14008eb0b: 48 8b 44 24 60              	movq	0x60(%rsp), %rax
14008eb10: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008eb14: e8 b7 b5 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008eb19: 48 83 7c 24 38 00           	cmpq	$0x0, 0x38(%rsp)
14008eb1f: 0f 84 2b 01 00 00           	je	0x14008ec50 <.text+0x8dc50>
14008eb25: 45 31 c0                    	xorl	%r8d, %r8d
14008eb28: ba 2f 00 00 00              	movl	$0x2f, %edx
14008eb2d: 48 8d 4c 24 30              	leaq	0x30(%rsp), %rcx
14008eb32: 48 8d 7c 24 30              	leaq	0x30(%rsp), %rdi
14008eb37: e8 84 89 02 00              	callq	0x1400b74c0 <.text+0xb64c0>
14008eb3c: 48 83 f8 ff                 	cmpq	$-0x1, %rax
14008eb40: 0f 84 6a 01 00 00           	je	0x14008ecb0 <.text+0x8dcb0>
14008eb46: 45 31 c0                    	xorl	%r8d, %r8d
14008eb49: 48 8b 4c 24 30              	movq	0x30(%rsp), %rcx
14008eb4e: 48 39 d9                    	cmpq	%rbx, %rcx
14008eb51: 0f 84 93 fe ff ff           	je	0x14008e9ea <.text+0x8d9ea>
14008eb57: 48 8b 44 24 40              	movq	0x40(%rsp), %rax
14008eb5c: 44 88 44 24 2f              	movb	%r8b, 0x2f(%rsp)
14008eb61: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008eb65: e8 66 b5 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008eb6a: 44 0f b6 44 24 2f           	movzbl	0x2f(%rsp), %r8d
14008eb70: 44 89 c0                    	movl	%r8d, %eax
14008eb73: 48 83 c4 78                 	addq	$0x78, %rsp
14008eb77: 5b                          	popq	%rbx
14008eb78: 5e                          	popq	%rsi
14008eb79: 5f                          	popq	%rdi
14008eb7a: 41 5e                       	popq	%r14
14008eb7c: c3                          	retq
14008eb7d: 0f 1f 00                    	nopl	(%rax)
14008eb80: 4c 8b 49 08                 	movq	0x8(%rcx), %r9
14008eb84: 49 83 f9 06                 	cmpq	$0x6, %r9
14008eb88: 0f 86 4e 9b 04 00           	jbe	0x1400d86dc <.text+0xd76dc>
14008eb8e: 48 8d 44 24 60              	leaq	0x60(%rsp), %rax
14008eb93: 4c 8d 74 24 50              	leaq	0x50(%rsp), %r14
14008eb98: 48 89 44 24 50              	movq	%rax, 0x50(%rsp)
14008eb9d: 48 89 c6                    	movq	%rax, %rsi
14008eba0: 48 8b 01                    	movq	(%rcx), %rax
14008eba3: 48 8d 4c 24 50              	leaq	0x50(%rsp), %rcx
14008eba8: 48 8d 7c 24 30              	leaq	0x30(%rsp), %rdi
14008ebad: 48 8d 50 07                 	leaq	0x7(%rax), %rdx
14008ebb1: 4e 8d 04 08                 	leaq	(%rax,%r9), %r8
14008ebb5: e8 e6 f8 ff ff              	callq	0x14008e4a0 <.text+0x8d4a0>
14008ebba: 48 8b 44 24 30              	movq	0x30(%rsp), %rax
14008ebbf: 49 89 c0                    	movq	%rax, %r8
14008ebc2: 48 39 d8                    	cmpq	%rbx, %rax
14008ebc5: 0f 84 a5 00 00 00           	je	0x14008ec70 <.text+0x8dc70>
14008ebcb: 48 8b 4c 24 50              	movq	0x50(%rsp), %rcx
14008ebd0: 48 8b 54 24 58              	movq	0x58(%rsp), %rdx
14008ebd5: 48 39 f1                    	cmpq	%rsi, %rcx
14008ebd8: 0f 85 f4 fe ff ff           	jne	0x14008ead2 <.text+0x8dad2>
14008ebde: 48 85 d2                    	testq	%rdx, %rdx
14008ebe1: 74 57                       	je	0x14008ec3a <.text+0x8dc3a>
14008ebe3: 48 83 fa 01                 	cmpq	$0x1, %rdx
14008ebe7: 0f 84 13 02 00 00           	je	0x14008ee00 <.text+0x8de00>
14008ebed: 89 d0                       	movl	%edx, %eax
14008ebef: 83 fa 40                    	cmpl	$0x40, %edx
14008ebf2: 0f 83 3e 02 00 00           	jae	0x14008ee36 <.text+0x8de36>
14008ebf8: f6 c2 20                    	testb	$0x20, %dl
14008ebfb: 0f 85 29 03 00 00           	jne	0x14008ef2a <.text+0x8df2a>
14008ec01: f6 c2 10                    	testb	$0x10, %dl
14008ec04: 0f 85 5d 03 00 00           	jne	0x14008ef67 <.text+0x8df67>
14008ec0a: f6 c2 08                    	testb	$0x8, %dl
14008ec0d: 0f 85 98 03 00 00           	jne	0x14008efab <.text+0x8dfab>
14008ec13: f6 c2 04                    	testb	$0x4, %dl
14008ec16: 0f 85 70 03 00 00           	jne	0x14008ef8c <.text+0x8df8c>
14008ec1c: 85 d2                       	testl	%edx, %edx
14008ec1e: 74 1a                       	je	0x14008ec3a <.text+0x8dc3a>
14008ec20: 0f b6 54 24 60              	movzbl	0x60(%rsp), %edx
14008ec25: 41 88 10                    	movb	%dl, (%r8)
14008ec28: a8 02                       	testb	$0x2, %al
14008ec2a: 0f 85 9c 03 00 00           	jne	0x14008efcc <.text+0x8dfcc>
14008ec30: 4c 8b 44 24 30              	movq	0x30(%rsp), %r8
14008ec35: 48 8b 54 24 58              	movq	0x58(%rsp), %rdx
14008ec3a: 41 c6 04 10 00              	movb	$0x0, (%r8,%rdx)
14008ec3f: 48 8b 44 24 50              	movq	0x50(%rsp), %rax
14008ec44: 48 89 54 24 38              	movq	%rdx, 0x38(%rsp)
14008ec49: e9 b0 fe ff ff              	jmp	0x14008eafe <.text+0x8dafe>
14008ec4e: 66 90                       	nop
14008ec50: 48 8b 4c 24 30              	movq	0x30(%rsp), %rcx
14008ec55: 45 31 c0                    	xorl	%r8d, %r8d
14008ec58: 48 39 d9                    	cmpq	%rbx, %rcx
14008ec5b: 0f 85 f6 fe ff ff           	jne	0x14008eb57 <.text+0x8db57>
14008ec61: e9 84 fd ff ff              	jmp	0x14008e9ea <.text+0x8d9ea>
14008ec66: 66 2e 0f 1f 84 00 00 00 00 00       	nopw	%cs:(%rax,%rax)
14008ec70: 48 8b 44 24 50              	movq	0x50(%rsp), %rax
14008ec75: 48 8b 54 24 58              	movq	0x58(%rsp), %rdx
14008ec7a: 48 39 f0                    	cmpq	%rsi, %rax
14008ec7d: 0f 84 5b ff ff ff           	je	0x14008ebde <.text+0x8dbde>
14008ec83: 48 89 44 24 30              	movq	%rax, 0x30(%rsp)
14008ec88: f3 0f 7e 44 24 58           	movq	0x58(%rsp), %xmm0
14008ec8e: 0f 16 44 24 60              	movhps	0x60(%rsp), %xmm0       # xmm0 = xmm0[0,1],mem[0,1]
14008ec93: 0f 11 44 24 38              	movups	%xmm0, 0x38(%rsp)
14008ec98: 48 8d 44 24 60              	leaq	0x60(%rsp), %rax
14008ec9d: 48 89 74 24 50              	movq	%rsi, 0x50(%rsp)
14008eca2: 48 89 c6                    	movq	%rax, %rsi
14008eca5: e9 54 fe ff ff              	jmp	0x14008eafe <.text+0x8dafe>
14008ecaa: 66 0f 1f 44 00 00           	nopw	(%rax,%rax)
14008ecb0: 45 31 c0                    	xorl	%r8d, %r8d
14008ecb3: ba 20 00 00 00              	movl	$0x20, %edx
14008ecb8: 48 89 f9                    	movq	%rdi, %rcx
14008ecbb: e8 00 88 02 00              	callq	0x1400b74c0 <.text+0xb64c0>
14008ecc0: 48 83 f8 ff                 	cmpq	$-0x1, %rax
14008ecc4: 0f 85 7c fe ff ff           	jne	0x14008eb46 <.text+0x8db46>
14008ecca: 49 c7 c0 ff ff ff ff        	movq	$-0x1, %r8
14008ecd1: ba 3a 00 00 00              	movl	$0x3a, %edx
14008ecd6: 48 89 f9                    	movq	%rdi, %rcx
14008ecd9: e8 72 89 02 00              	callq	0x1400b7650 <.text+0xb6650>
14008ecde: 48 83 f8 ff                 	cmpq	$-0x1, %rax
14008ece2: 0f 84 a0 00 00 00           	je	0x14008ed88 <.text+0x8dd88>
14008ece8: 4c 8b 4c 24 38              	movq	0x38(%rsp), %r9
14008eced: 4c 8d 40 01                 	leaq	0x1(%rax), %r8
14008ecf1: 4d 39 c1                    	cmpq	%r8, %r9
14008ecf4: 0f 82 cf 99 04 00           	jb	0x1400d86c9 <.text+0xd76c9>
14008ecfa: 48 8b 44 24 30              	movq	0x30(%rsp), %rax
14008ecff: 4c 89 f1                    	movq	%r14, %rcx
14008ed02: 48 89 74 24 50              	movq	%rsi, 0x50(%rsp)
14008ed07: 4a 8d 14 00                 	leaq	(%rax,%r8), %rdx
14008ed0b: 4e 8d 04 08                 	leaq	(%rax,%r9), %r8
14008ed0f: e8 8c f7 ff ff              	callq	0x14008e4a0 <.text+0x8d4a0>
14008ed14: 48 8b 4c 24 58              	movq	0x58(%rsp), %rcx
14008ed19: 4c 8b 44 24 50              	movq	0x50(%rsp), %r8
14008ed1e: 48 85 c9                    	testq	%rcx, %rcx
14008ed21: 0f 84 f0 00 00 00           	je	0x14008ee17 <.text+0x8de17>
14008ed27: 48 83 f9 05                 	cmpq	$0x5, %rcx
14008ed2b: 0f 87 e6 00 00 00           	ja	0x14008ee17 <.text+0x8de17>
14008ed31: 4c 01 c1                    	addq	%r8, %rcx
14008ed34: 4c 89 c2                    	movq	%r8, %rdx
14008ed37: 66 0f 1f 84 00 00 00 00 00  	nopw	(%rax,%rax)
14008ed40: 0f b6 02                    	movzbl	(%rdx), %eax
14008ed43: 83 e8 30                    	subl	$0x30, %eax
14008ed46: 3c 09                       	cmpb	$0x9, %al
14008ed48: 0f 87 c9 00 00 00           	ja	0x14008ee17 <.text+0x8de17>
14008ed4e: 48 83 c2 01                 	addq	$0x1, %rdx
14008ed52: 48 39 d1                    	cmpq	%rdx, %rcx
14008ed55: 75 e9                       	jne	0x14008ed40 <.text+0x8dd40>
14008ed57: 4c 89 c1                    	movq	%r8, %rcx
14008ed5a: e8 41 6b 01 00              	callq	0x1400a58a0 <.text+0xa48a0>
14008ed5f: 4c 8b 44 24 50              	movq	0x50(%rsp), %r8
14008ed64: 83 e8 01                    	subl	$0x1, %eax
14008ed67: 3d fe ff 00 00              	cmpl	$0xfffe, %eax           # imm = 0xFFFE
14008ed6c: 0f 87 a5 00 00 00           	ja	0x14008ee17 <.text+0x8de17>
14008ed72: 49 39 f0                    	cmpq	%rsi, %r8
14008ed75: 74 11                       	je	0x14008ed88 <.text+0x8dd88>
14008ed77: 48 8b 44 24 60              	movq	0x60(%rsp), %rax
14008ed7c: 4c 89 c1                    	movq	%r8, %rcx
14008ed7f: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008ed83: e8 48 b3 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008ed88: 41 b8 01 00 00 00           	movl	$0x1, %r8d
14008ed8e: e9 b6 fd ff ff              	jmp	0x14008eb49 <.text+0x8db49>
14008ed93: 0f 1f 44 00 00              	nopl	(%rax,%rax)
14008ed98: 48 8b 44 24 50              	movq	0x50(%rsp), %rax
14008ed9d: 48 8b 54 24 58              	movq	0x58(%rsp), %rdx
14008eda2: 48 39 f0                    	cmpq	%rsi, %rax
14008eda5: 0f 85 d8 fe ff ff           	jne	0x14008ec83 <.text+0x8dc83>
14008edab: 0f 1f 44 00 00              	nopl	(%rax,%rax)
14008edb0: 48 85 d2                    	testq	%rdx, %rdx
14008edb3: 0f 84 81 fe ff ff           	je	0x14008ec3a <.text+0x8dc3a>
14008edb9: 48 83 fa 01                 	cmpq	$0x1, %rdx
14008edbd: 74 41                       	je	0x14008ee00 <.text+0x8de00>
14008edbf: 89 d0                       	movl	%edx, %eax
14008edc1: 83 fa 40                    	cmpl	$0x40, %edx
14008edc4: 0f 83 e6 00 00 00           	jae	0x14008eeb0 <.text+0x8deb0>
14008edca: f6 c2 20                    	testb	$0x20, %dl
14008edcd: 0f 85 57 01 00 00           	jne	0x14008ef2a <.text+0x8df2a>
14008edd3: f6 c2 10                    	testb	$0x10, %dl
14008edd6: 0f 85 8b 01 00 00           	jne	0x14008ef67 <.text+0x8df67>
14008eddc: a8 08                       	testb	$0x8, %al
14008edde: 0f 85 c7 01 00 00           	jne	0x14008efab <.text+0x8dfab>
14008ede4: a8 04                       	testb	$0x4, %al
14008ede6: 0f 85 a0 01 00 00           	jne	0x14008ef8c <.text+0x8df8c>
14008edec: 85 c0                       	testl	%eax, %eax
14008edee: 0f 84 46 fe ff ff           	je	0x14008ec3a <.text+0x8dc3a>
14008edf4: e9 27 fe ff ff              	jmp	0x14008ec20 <.text+0x8dc20>
14008edf9: 0f 1f 80 00 00 00 00        	nopl	(%rax)
14008ee00: 0f b6 44 24 60              	movzbl	0x60(%rsp), %eax
14008ee05: 41 88 00                    	movb	%al, (%r8)
14008ee08: 4c 8b 44 24 30              	movq	0x30(%rsp), %r8
14008ee0d: 48 8b 54 24 58              	movq	0x58(%rsp), %rdx
14008ee12: e9 23 fe ff ff              	jmp	0x14008ec3a <.text+0x8dc3a>
14008ee17: 49 39 f0                    	cmpq	%rsi, %r8
14008ee1a: 0f 84 26 fd ff ff           	je	0x14008eb46 <.text+0x8db46>
14008ee20: 48 8b 44 24 60              	movq	0x60(%rsp), %rax
14008ee25: 4c 89 c1                    	movq	%r8, %rcx
14008ee28: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008ee2c: e8 9f b2 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008ee31: e9 10 fd ff ff              	jmp	0x14008eb46 <.text+0x8db46>
14008ee36: 89 d0                       	movl	%edx, %eax
14008ee38: 49 8d 0c 00                 	leaq	(%r8,%rax), %rcx
14008ee3c: 48 01 f0                    	addq	%rsi, %rax
14008ee3f: f3 0f 6f 40 c0              	movdqu	-0x40(%rax), %xmm0
14008ee44: 0f 11 41 c0                 	movups	%xmm0, -0x40(%rcx)
14008ee48: f3 0f 6f 40 d0              	movdqu	-0x30(%rax), %xmm0
14008ee4d: 0f 11 41 d0                 	movups	%xmm0, -0x30(%rcx)
14008ee51: f3 0f 6f 40 e0              	movdqu	-0x20(%rax), %xmm0
14008ee56: 0f 11 41 e0                 	movups	%xmm0, -0x20(%rcx)
14008ee5a: f3 0f 6f 40 f0              	movdqu	-0x10(%rax), %xmm0
14008ee5f: 8d 42 ff                    	leal	-0x1(%rdx), %eax
14008ee62: 0f 11 41 f0                 	movups	%xmm0, -0x10(%rcx)
14008ee66: 83 f8 40                    	cmpl	$0x40, %eax
14008ee69: 0f 82 c1 fd ff ff           	jb	0x14008ec30 <.text+0x8dc30>
14008ee6f: 83 e0 c0                    	andl	$-0x40, %eax
14008ee72: 31 c9                       	xorl	%ecx, %ecx
14008ee74: 89 ca                       	movl	%ecx, %edx
14008ee76: 83 c1 40                    	addl	$0x40, %ecx
14008ee79: f3 0f 6f 1c 16              	movdqu	(%rsi,%rdx), %xmm3
14008ee7e: f3 0f 6f 54 16 10           	movdqu	0x10(%rsi,%rdx), %xmm2
14008ee84: f3 0f 6f 4c 16 20           	movdqu	0x20(%rsi,%rdx), %xmm1
14008ee8a: f3 0f 6f 44 16 30           	movdqu	0x30(%rsi,%rdx), %xmm0
14008ee90: 41 0f 11 1c 10              	movups	%xmm3, (%r8,%rdx)
14008ee95: 41 0f 11 54 10 10           	movups	%xmm2, 0x10(%r8,%rdx)
14008ee9b: 41 0f 11 4c 10 20           	movups	%xmm1, 0x20(%r8,%rdx)
14008eea1: 41 0f 11 44 10 30           	movups	%xmm0, 0x30(%r8,%rdx)
14008eea7: 39 c1                       	cmpl	%eax, %ecx
14008eea9: 72 c9                       	jb	0x14008ee74 <.text+0x8de74>
14008eeab: e9 80 fd ff ff              	jmp	0x14008ec30 <.text+0x8dc30>
14008eeb0: 89 d0                       	movl	%edx, %eax
14008eeb2: 49 8d 0c 00                 	leaq	(%r8,%rax), %rcx
14008eeb6: 48 01 f0                    	addq	%rsi, %rax
14008eeb9: f3 0f 6f 40 c0              	movdqu	-0x40(%rax), %xmm0
14008eebe: 0f 11 41 c0                 	movups	%xmm0, -0x40(%rcx)
14008eec2: f3 0f 6f 40 d0              	movdqu	-0x30(%rax), %xmm0
14008eec7: 0f 11 41 d0                 	movups	%xmm0, -0x30(%rcx)
14008eecb: f3 0f 6f 40 e0              	movdqu	-0x20(%rax), %xmm0
14008eed0: 0f 11 41 e0                 	movups	%xmm0, -0x20(%rcx)
14008eed4: f3 0f 6f 40 f0              	movdqu	-0x10(%rax), %xmm0
14008eed9: 8d 42 ff                    	leal	-0x1(%rdx), %eax
14008eedc: 0f 11 41 f0                 	movups	%xmm0, -0x10(%rcx)
14008eee0: 83 f8 40                    	cmpl	$0x40, %eax
14008eee3: 0f 82 47 fd ff ff           	jb	0x14008ec30 <.text+0x8dc30>
14008eee9: 83 e0 c0                    	andl	$-0x40, %eax
14008eeec: 31 c9                       	xorl	%ecx, %ecx
14008eeee: 89 ca                       	movl	%ecx, %edx
14008eef0: 83 c1 40                    	addl	$0x40, %ecx
14008eef3: f3 0f 6f 1c 16              	movdqu	(%rsi,%rdx), %xmm3
14008eef8: f3 0f 6f 54 16 10           	movdqu	0x10(%rsi,%rdx), %xmm2
14008eefe: f3 0f 6f 4c 16 20           	movdqu	0x20(%rsi,%rdx), %xmm1
14008ef04: f3 0f 6f 44 16 30           	movdqu	0x30(%rsi,%rdx), %xmm0
14008ef0a: 41 0f 11 1c 10              	movups	%xmm3, (%r8,%rdx)
14008ef0f: 41 0f 11 54 10 10           	movups	%xmm2, 0x10(%r8,%rdx)
14008ef15: 41 0f 11 4c 10 20           	movups	%xmm1, 0x20(%r8,%rdx)
14008ef1b: 41 0f 11 44 10 30           	movups	%xmm0, 0x30(%r8,%rdx)
14008ef21: 39 c1                       	cmpl	%eax, %ecx
14008ef23: 72 c9                       	jb	0x14008eeee <.text+0x8deee>
14008ef25: e9 06 fd ff ff              	jmp	0x14008ec30 <.text+0x8dc30>
14008ef2a: f3 0f 6f 44 24 60           	movdqu	0x60(%rsp), %xmm0
14008ef30: 49 8d 54 00 20              	leaq	0x20(%r8,%rax), %rdx
14008ef35: 41 0f 11 00                 	movups	%xmm0, (%r8)
14008ef39: f3 0f 6f 44 24 70           	movdqu	0x70(%rsp), %xmm0
14008ef3f: 41 0f 11 40 10              	movups	%xmm0, 0x10(%r8)
14008ef44: f3 0f 6f 44 06 e0           	movdqu	-0x20(%rsi,%rax), %xmm0
14008ef4a: 0f 11 42 c0                 	movups	%xmm0, -0x40(%rdx)
14008ef4e: f3 0f 6f 44 06 f0           	movdqu	-0x10(%rsi,%rax), %xmm0
14008ef54: 0f 11 42 d0                 	movups	%xmm0, -0x30(%rdx)
14008ef58: 4c 8b 44 24 30              	movq	0x30(%rsp), %r8
14008ef5d: 48 8b 54 24 58              	movq	0x58(%rsp), %rdx
14008ef62: e9 d3 fc ff ff              	jmp	0x14008ec3a <.text+0x8dc3a>
14008ef67: f3 0f 6f 44 24 60           	movdqu	0x60(%rsp), %xmm0
14008ef6d: 41 0f 11 00                 	movups	%xmm0, (%r8)
14008ef71: f3 0f 6f 44 06 f0           	movdqu	-0x10(%rsi,%rax), %xmm0
14008ef77: 41 0f 11 44 00 f0           	movups	%xmm0, -0x10(%r8,%rax)
14008ef7d: 4c 8b 44 24 30              	movq	0x30(%rsp), %r8
14008ef82: 48 8b 54 24 58              	movq	0x58(%rsp), %rdx
14008ef87: e9 ae fc ff ff              	jmp	0x14008ec3a <.text+0x8dc3a>
14008ef8c: 8b 54 24 60                 	movl	0x60(%rsp), %edx
14008ef90: 41 89 10                    	movl	%edx, (%r8)
14008ef93: 8b 54 06 fc                 	movl	-0x4(%rsi,%rax), %edx
14008ef97: 41 89 54 00 fc              	movl	%edx, -0x4(%r8,%rax)
14008ef9c: 4c 8b 44 24 30              	movq	0x30(%rsp), %r8
14008efa1: 48 8b 54 24 58              	movq	0x58(%rsp), %rdx
14008efa6: e9 8f fc ff ff              	jmp	0x14008ec3a <.text+0x8dc3a>
14008efab: 48 8b 54 24 60              	movq	0x60(%rsp), %rdx
14008efb0: 49 89 10                    	movq	%rdx, (%r8)
14008efb3: 48 8b 54 06 f8              	movq	-0x8(%rsi,%rax), %rdx
14008efb8: 49 89 54 00 f8              	movq	%rdx, -0x8(%r8,%rax)
14008efbd: 4c 8b 44 24 30              	movq	0x30(%rsp), %r8
14008efc2: 48 8b 54 24 58              	movq	0x58(%rsp), %rdx
14008efc7: e9 6e fc ff ff              	jmp	0x14008ec3a <.text+0x8dc3a>
14008efcc: 0f b7 54 06 fe              	movzwl	-0x2(%rsi,%rax), %edx
14008efd1: 66 41 89 54 00 fe           	movw	%dx, -0x2(%r8,%rax)
14008efd7: 4c 8b 44 24 30              	movq	0x30(%rsp), %r8
14008efdc: 48 8b 54 24 58              	movq	0x58(%rsp), %rdx
14008efe1: e9 54 fc ff ff              	jmp	0x14008ec3a <.text+0x8dc3a>
14008efe6: e9 0f 97 04 00              	jmp	0x1400d86fa <.text+0xd76fa>
14008efeb: 0f 1f 44 00 00              	nopl	(%rax,%rax)
14008eff0: 41 57                       	pushq	%r15
14008eff2: 41 56                       	pushq	%r14
14008eff4: 41 55                       	pushq	%r13
14008eff6: 41 54                       	pushq	%r12
14008eff8: 55                          	pushq	%rbp
14008eff9: 57                          	pushq	%rdi
14008effa: 56                          	pushq	%rsi
14008effb: 53                          	pushq	%rbx
14008effc: 48 81 ec a8 00 00 00        	subq	$0xa8, %rsp
14008f003: 89 cf                       	movl	%ecx, %edi
14008f005: 49 89 d5                    	movq	%rdx, %r13
14008f008: 4c 89 c6                    	movq	%r8, %rsi
14008f00b: 4c 89 cd                    	movq	%r9, %rbp
14008f00e: 45 31 c0                    	xorl	%r8d, %r8d
14008f011: 45 31 c9                    	xorl	%r9d, %r9d
14008f014: 48 8d 8c 24 80 00 00 00     	leaq	0x80(%rsp), %rcx
14008f01c: 48 c7 84 24 80 00 00 00 00 00 00 00 	movq	$0x0, 0x80(%rsp)
14008f028: 48 8d 15 6b c6 06 00        	leaq	0x6c66b(%rip), %rdx     # 0x1400fb69a
14008f02f: 48 c7 84 24 88 00 00 00 00 00 00 00 	movq	$0x0, 0x88(%rsp)
14008f03b: e8 20 1c 00 00              	callq	0x140090c60 <.text+0x8fc60>
14008f040: 85 c0                       	testl	%eax, %eax
14008f042: 0f 88 50 03 00 00           	js	0x14008f398 <.text+0x8e398>
14008f048: c7 44 24 20 00 00 00 00     	movl	$0x0, 0x20(%rsp)
14008f050: 48 8b 8c 24 80 00 00 00     	movq	0x80(%rsp), %rcx
14008f058: 41 b9 20 00 00 00           	movl	$0x20, %r9d
14008f05e: 4c 8d 05 43 c6 06 00        	leaq	0x6c643(%rip), %r8      # 0x1400fb6a8
14008f065: 48 8d 15 5c c6 06 00        	leaq	0x6c65c(%rip), %rdx     # 0x1400fb6c8
14008f06c: e8 df 1b 00 00              	callq	0x140090c50 <.text+0x8fc50>
14008f071: c7 44 24 74 00 00 00 00     	movl	$0x0, 0x74(%rsp)
14008f079: c7 44 24 78 00 00 00 00     	movl	$0x0, 0x78(%rsp)
14008f081: 41 89 c6                    	movl	%eax, %r14d
14008f084: 85 c0                       	testl	%eax, %eax
14008f086: 0f 89 24 02 00 00           	jns	0x14008f2b0 <.text+0x8e2b0>
14008f08c: b9 40 00 00 00              	movl	$0x40, %ecx
14008f091: 41 bc 40 00 00 00           	movl	$0x40, %r12d
14008f097: e8 44 b0 03 00              	callq	0x1400ca0e0 <.text+0xc90e0>
14008f09c: 66 0f ef c0                 	pxor	%xmm0, %xmm0
14008f0a0: 41 b9 40 00 00 00           	movl	$0x40, %r9d
14008f0a6: 48 89 c3                    	movq	%rax, %rbx
14008f0a9: 48 8d 40 40                 	leaq	0x40(%rax), %rax
14008f0ad: 0f 11 43 21                 	movups	%xmm0, 0x21(%rbx)
14008f0b1: 48 89 44 24 50              	movq	%rax, 0x50(%rsp)
14008f0b6: c6 03 00                    	movb	$0x0, (%rbx)
14008f0b9: 0f 11 43 01                 	movups	%xmm0, 0x1(%rbx)
14008f0bd: 0f 11 43 11                 	movups	%xmm0, 0x11(%rbx)
14008f0c1: 0f 11 43 30                 	movups	%xmm0, 0x30(%rbx)
14008f0c5: 45 85 f6                    	testl	%r14d, %r14d
14008f0c8: 0f 88 61 02 00 00           	js	0x14008f32f <.text+0x8e32f>
14008f0ce: 4c 89 6c 24 20              	movq	%r13, 0x20(%rsp)
14008f0d3: 48 8b 8c 24 80 00 00 00     	movq	0x80(%rsp), %rcx
14008f0db: 49 89 d8                    	movq	%rbx, %r8
14008f0de: 48 8d 94 24 88 00 00 00     	leaq	0x88(%rsp), %rdx
14008f0e6: c7 44 24 30 00 00 00 00     	movl	$0x0, 0x30(%rsp)
14008f0ee: c7 44 24 28 20 00 00 00     	movl	$0x20, 0x28(%rsp)
14008f0f6: e8 85 1b 00 00              	callq	0x140090c80 <.text+0x8fc80>
14008f0fb: f3 0f 6f 06                 	movdqu	(%rsi), %xmm0
14008f0ff: c7 44 24 7c 00 00 00 00     	movl	$0x0, 0x7c(%rsp)
14008f107: 41 89 c6                    	movl	%eax, %r14d
14008f10a: 48 8b 8c 24 88 00 00 00     	movq	0x88(%rsp), %rcx
14008f112: 0f 29 84 24 90 00 00 00     	movaps	%xmm0, 0x90(%rsp)
14008f11a: 85 c0                       	testl	%eax, %eax
14008f11c: 0f 88 29 02 00 00           	js	0x14008f34b <.text+0x8e34b>
14008f122: 48 8d 44 24 7c              	leaq	0x7c(%rsp), %rax
14008f127: 4c 8d ac 24 90 00 00 00     	leaq	0x90(%rsp), %r13
14008f12f: c7 44 24 48 01 00 00 00     	movl	$0x1, 0x48(%rsp)
14008f137: 44 8b 84 24 10 01 00 00     	movl	0x110(%rsp), %r8d
14008f13f: 48 89 44 24 58              	movq	%rax, 0x58(%rsp)
14008f144: 45 31 c9                    	xorl	%r9d, %r9d
14008f147: 48 89 ea                    	movq	%rbp, %rdx
14008f14a: 48 89 44 24 40              	movq	%rax, 0x40(%rsp)
14008f14f: c7 44 24 38 00 00 00 00     	movl	$0x0, 0x38(%rsp)
14008f157: 48 c7 44 24 30 00 00 00 00  	movq	$0x0, 0x30(%rsp)
14008f160: c7 44 24 28 10 00 00 00     	movl	$0x10, 0x28(%rsp)
14008f168: 4c 89 6c 24 20              	movq	%r13, 0x20(%rsp)
14008f16d: 85 ff                       	testl	%edi, %edi
14008f16f: 0f 85 2b 02 00 00           	jne	0x14008f3a0 <.text+0x8e3a0>
14008f175: e8 5e 1b 00 00              	callq	0x140090cd8 <.text+0x8fcd8>
14008f17a: 48 8b 8c 24 88 00 00 00     	movq	0x88(%rsp), %rcx
14008f182: 41 89 c6                    	movl	%eax, %r14d
14008f185: 85 c0                       	testl	%eax, %eax
14008f187: 0f 88 be 01 00 00           	js	0x14008f34b <.text+0x8e34b>
14008f18d: 48 8b 8c 24 18 01 00 00     	movq	0x118(%rsp), %rcx
14008f195: 48 8b 84 24 18 01 00 00     	movq	0x118(%rsp), %rax
14008f19d: 44 8b 4c 24 7c              	movl	0x7c(%rsp), %r9d
14008f1a2: 48 8b 49 10                 	movq	0x10(%rcx), %rcx
14008f1a6: 4c 8b 18                    	movq	(%rax), %r11
14008f1a9: 4d 89 ce                    	movq	%r9, %r14
14008f1ac: 44 89 ca                    	movl	%r9d, %edx
14008f1af: 49 89 cf                    	movq	%rcx, %r15
14008f1b2: 48 89 4c 24 60              	movq	%rcx, 0x60(%rsp)
14008f1b7: 4c 89 d8                    	movq	%r11, %rax
14008f1ba: 4d 29 df                    	subq	%r11, %r15
14008f1bd: 4d 39 cf                    	cmpq	%r9, %r15
14008f1c0: 0f 83 ea 01 00 00           	jae	0x14008f3b0 <.text+0x8e3b0>
14008f1c6: 4c 89 c9                    	movq	%r9, %rcx
14008f1c9: 4c 89 5c 24 68              	movq	%r11, 0x68(%rsp)
14008f1ce: 4c 89 4c 24 60              	movq	%r9, 0x60(%rsp)
14008f1d3: e8 08 af 03 00              	callq	0x1400ca0e0 <.text+0xc90e0>
14008f1d8: 4c 8b 4c 24 60              	movq	0x60(%rsp), %r9
14008f1dd: 48 89 c1                    	movq	%rax, %rcx
14008f1e0: 31 d2                       	xorl	%edx, %edx
14008f1e2: 4e 8d 14 08                 	leaq	(%rax,%r9), %r10
14008f1e6: 4d 89 c8                    	movq	%r9, %r8
14008f1e9: 4c 89 54 24 60              	movq	%r10, 0x60(%rsp)
14008f1ee: e8 5d 64 01 00              	callq	0x1400a5650 <.text+0xa4650>
14008f1f3: 4c 8b 54 24 60              	movq	0x60(%rsp), %r10
14008f1f8: 4c 8b 5c 24 68              	movq	0x68(%rsp), %r11
14008f1fd: 48 89 c1                    	movq	%rax, %rcx
14008f200: 48 8b 84 24 18 01 00 00     	movq	0x118(%rsp), %rax
14008f208: 48 89 08                    	movq	%rcx, (%rax)
14008f20b: 4c 89 50 08                 	movq	%r10, 0x8(%rax)
14008f20f: 4c 89 50 10                 	movq	%r10, 0x10(%rax)
14008f213: 4d 85 db                    	testq	%r11, %r11
14008f216: 74 1b                       	je	0x14008f233 <.text+0x8e233>
14008f218: 4c 89 d9                    	movq	%r11, %rcx
14008f21b: 4c 89 fa                    	movq	%r15, %rdx
14008f21e: e8 ad ae 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008f223: 48 8b 84 24 18 01 00 00     	movq	0x118(%rsp), %rax
14008f22b: 44 8b 74 24 7c              	movl	0x7c(%rsp), %r14d
14008f230: 48 8b 08                    	movq	(%rax), %rcx
14008f233: f3 0f 6f 06                 	movdqu	(%rsi), %xmm0
14008f237: 85 ff                       	testl	%edi, %edi
14008f239: 48 8b 7c 24 58              	movq	0x58(%rsp), %rdi
14008f23e: 48 89 4c 24 30              	movq	%rcx, 0x30(%rsp)
14008f243: 48 8b 84 24 88 00 00 00     	movq	0x88(%rsp), %rax
14008f24b: 44 89 74 24 38              	movl	%r14d, 0x38(%rsp)
14008f250: 48 89 ea                    	movq	%rbp, %rdx
14008f253: 41 b9 00 00 00 00           	movl	$0x0, %r9d
14008f259: c7 44 24 48 01 00 00 00     	movl	$0x1, 0x48(%rsp)
14008f261: 44 8b 84 24 10 01 00 00     	movl	0x110(%rsp), %r8d
14008f269: 48 89 7c 24 40              	movq	%rdi, 0x40(%rsp)
14008f26e: 48 89 c1                    	movq	%rax, %rcx
14008f271: c7 44 24 28 10 00 00 00     	movl	$0x10, 0x28(%rsp)
14008f279: 4c 89 6c 24 20              	movq	%r13, 0x20(%rsp)
14008f27e: 0f 29 84 24 90 00 00 00     	movaps	%xmm0, 0x90(%rsp)
14008f286: 0f 84 7c 01 00 00           	je	0x14008f408 <.text+0x8e408>
14008f28c: e8 1f 1a 00 00              	callq	0x140090cb0 <.text+0x8fcb0>
14008f291: 41 89 c6                    	movl	%eax, %r14d
14008f294: 85 c0                       	testl	%eax, %eax
14008f296: 0f 89 7c 01 00 00           	jns	0x14008f418 <.text+0x8e418>
14008f29c: 48 8b 8c 24 88 00 00 00     	movq	0x88(%rsp), %rcx
14008f2a4: e9 a2 00 00 00              	jmp	0x14008f34b <.text+0x8e34b>
14008f2a9: 0f 1f 80 00 00 00 00        	nopl	(%rax)
14008f2b0: 48 8d 44 24 78              	leaq	0x78(%rsp), %rax
14008f2b5: c7 44 24 28 00 00 00 00     	movl	$0x0, 0x28(%rsp)
14008f2bd: 4c 8d 44 24 74              	leaq	0x74(%rsp), %r8
14008f2c2: 48 8b 8c 24 80 00 00 00     	movq	0x80(%rsp), %rcx
14008f2ca: 48 89 44 24 20              	movq	%rax, 0x20(%rsp)
14008f2cf: 41 b9 04 00 00 00           	movl	$0x4, %r9d
14008f2d5: 48 8d 15 06 c4 06 00        	leaq	0x6c406(%rip), %rdx     # 0x1400fb6e2
14008f2dc: e8 97 19 00 00              	callq	0x140090c78 <.text+0x8fc78>
14008f2e1: 44 8b 64 24 74              	movl	0x74(%rsp), %r12d
14008f2e6: 41 89 c6                    	movl	%eax, %r14d
14008f2e9: 45 85 e4                    	testl	%r12d, %r12d
14008f2ec: 0f 84 9a fd ff ff           	je	0x14008f08c <.text+0x8e08c>
14008f2f2: 45 89 e1                    	movl	%r12d, %r9d
14008f2f5: 4c 89 c9                    	movq	%r9, %rcx
14008f2f8: 4c 89 4c 24 50              	movq	%r9, 0x50(%rsp)
14008f2fd: e8 de ad 03 00              	callq	0x1400ca0e0 <.text+0xc90e0>
14008f302: 4c 8b 4c 24 50              	movq	0x50(%rsp), %r9
14008f307: 48 89 c3                    	movq	%rax, %rbx
14008f30a: 4a 8d 04 08                 	leaq	(%rax,%r9), %rax
14008f30e: c6 03 00                    	movb	$0x0, (%rbx)
14008f311: 48 89 44 24 50              	movq	%rax, 0x50(%rsp)
14008f316: 41 83 fc 01                 	cmpl	$0x1, %r12d
14008f31a: 0f 85 50 01 00 00           	jne	0x14008f470 <.text+0x8e470>
14008f320: 41 b9 01 00 00 00           	movl	$0x1, %r9d
14008f326: 45 85 f6                    	testl	%r14d, %r14d
14008f329: 0f 89 9f fd ff ff           	jns	0x14008f0ce <.text+0x8e0ce>
14008f32f: f3 0f 6f 06                 	movdqu	(%rsi), %xmm0
14008f333: c7 44 24 7c 00 00 00 00     	movl	$0x0, 0x7c(%rsp)
14008f33b: 48 8b 8c 24 88 00 00 00     	movq	0x88(%rsp), %rcx
14008f343: 0f 29 84 24 90 00 00 00     	movaps	%xmm0, 0x90(%rsp)
14008f34b: 48 85 c9                    	testq	%rcx, %rcx
14008f34e: 74 05                       	je	0x14008f355 <.text+0x8e355>
14008f350: e8 6b 19 00 00              	callq	0x140090cc0 <.text+0x8fcc0>
14008f355: 48 8b 8c 24 80 00 00 00     	movq	0x80(%rsp), %rcx
14008f35d: 48 85 c9                    	testq	%rcx, %rcx
14008f360: 74 07                       	je	0x14008f369 <.text+0x8e369>
14008f362: 31 d2                       	xorl	%edx, %edx
14008f364: e8 7f 19 00 00              	callq	0x140090ce8 <.text+0x8fce8>
14008f369: 41 f7 d6                    	notl	%r14d
14008f36c: 44 89 e2                    	movl	%r12d, %edx
14008f36f: 48 89 d9                    	movq	%rbx, %rcx
14008f372: 41 c1 ee 1f                 	shrl	$0x1f, %r14d
14008f376: e8 55 ad 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008f37b: 44 89 f0                    	movl	%r14d, %eax
14008f37e: 48 81 c4 a8 00 00 00        	addq	$0xa8, %rsp
14008f385: 5b                          	popq	%rbx
14008f386: 5e                          	popq	%rsi
14008f387: 5f                          	popq	%rdi
14008f388: 5d                          	popq	%rbp
14008f389: 41 5c                       	popq	%r12
14008f38b: 41 5d                       	popq	%r13
14008f38d: 41 5e                       	popq	%r14
14008f38f: 41 5f                       	popq	%r15
14008f391: c3                          	retq
14008f392: 66 0f 1f 44 00 00           	nopw	(%rax,%rax)
14008f398: 45 31 f6                    	xorl	%r14d, %r14d
14008f39b: eb de                       	jmp	0x14008f37b <.text+0x8e37b>
14008f39d: 0f 1f 00                    	nopl	(%rax)
14008f3a0: e8 0b 19 00 00              	callq	0x140090cb0 <.text+0x8fcb0>
14008f3a5: e9 d0 fd ff ff              	jmp	0x14008f17a <.text+0x8e17a>
14008f3aa: 66 0f 1f 44 00 00           	nopw	(%rax,%rax)
14008f3b0: 48 8b 8c 24 18 01 00 00     	movq	0x118(%rsp), %rcx
14008f3b8: 48 8b 49 08                 	movq	0x8(%rcx), %rcx
14008f3bc: 49 89 c8                    	movq	%rcx, %r8
14008f3bf: 4d 29 d8                    	subq	%r11, %r8
14008f3c2: 4d 39 c8                    	cmpq	%r9, %r8
14008f3c5: 0f 83 1a 01 00 00           	jae	0x14008f4e5 <.text+0x8e4e5>
14008f3cb: 4d 85 c0                    	testq	%r8, %r8
14008f3ce: 0f 85 cc 00 00 00           	jne	0x14008f4a0 <.text+0x8e4a0>
14008f3d4: 4d 89 ce                    	movq	%r9, %r14
14008f3d7: 31 d2                       	xorl	%edx, %edx
14008f3d9: 4d 29 c6                    	subq	%r8, %r14
14008f3dc: 4d 89 f0                    	movq	%r14, %r8
14008f3df: e8 6c 62 01 00              	callq	0x1400a5650 <.text+0xa4650>
14008f3e4: 4e 8d 0c 30                 	leaq	(%rax,%r14), %r9
14008f3e8: 48 8b 84 24 18 01 00 00     	movq	0x118(%rsp), %rax
14008f3f0: 44 8b 74 24 7c              	movl	0x7c(%rsp), %r14d
14008f3f5: 4c 89 48 08                 	movq	%r9, 0x8(%rax)
14008f3f9: 48 8b 08                    	movq	(%rax), %rcx
14008f3fc: e9 32 fe ff ff              	jmp	0x14008f233 <.text+0x8e233>
14008f401: 0f 1f 80 00 00 00 00        	nopl	(%rax)
14008f408: e8 cb 18 00 00              	callq	0x140090cd8 <.text+0x8fcd8>
14008f40d: e9 7f fe ff ff              	jmp	0x14008f291 <.text+0x8e291>
14008f412: 66 0f 1f 44 00 00           	nopw	(%rax,%rax)
14008f418: 48 8b 84 24 18 01 00 00     	movq	0x118(%rsp), %rax
14008f420: 8b 54 24 7c                 	movl	0x7c(%rsp), %edx
14008f424: 4c 8b 48 08                 	movq	0x8(%rax), %r9
14008f428: 48 8b 00                    	movq	(%rax), %rax
14008f42b: 4d 89 c8                    	movq	%r9, %r8
14008f42e: 49 29 c0                    	subq	%rax, %r8
14008f431: 49 39 d0                    	cmpq	%rdx, %r8
14008f434: 0f 82 96 00 00 00           	jb	0x14008f4d0 <.text+0x8e4d0>
14008f43a: 48 8b 8c 24 88 00 00 00     	movq	0x88(%rsp), %rcx
14008f442: 4c 39 c2                    	cmpq	%r8, %rdx
14008f445: 0f 83 00 ff ff ff           	jae	0x14008f34b <.text+0x8e34b>
14008f44b: 48 01 d0                    	addq	%rdx, %rax
14008f44e: 49 39 c1                    	cmpq	%rax, %r9
14008f451: 0f 84 f4 fe ff ff           	je	0x14008f34b <.text+0x8e34b>
14008f457: 48 8b b4 24 18 01 00 00     	movq	0x118(%rsp), %rsi
14008f45f: 48 89 46 08                 	movq	%rax, 0x8(%rsi)
14008f463: e9 e3 fe ff ff              	jmp	0x14008f34b <.text+0x8e34b>
14008f468: 0f 1f 84 00 00 00 00 00     	nopl	(%rax,%rax)
14008f470: 4d 8d 41 ff                 	leaq	-0x1(%r9), %r8
14008f474: 48 8d 4b 01                 	leaq	0x1(%rbx), %rcx
14008f478: 31 d2                       	xorl	%edx, %edx
14008f47a: 4c 89 4c 24 58              	movq	%r9, 0x58(%rsp)
14008f47f: e8 cc 61 01 00              	callq	0x1400a5650 <.text+0xa4650>
14008f484: 4c 8b 4c 24 58              	movq	0x58(%rsp), %r9
14008f489: 45 85 f6                    	testl	%r14d, %r14d
14008f48c: 0f 89 3c fc ff ff           	jns	0x14008f0ce <.text+0x8e0ce>
14008f492: e9 98 fe ff ff              	jmp	0x14008f32f <.text+0x8e32f>
14008f497: 66 0f 1f 84 00 00 00 00 00  	nopw	(%rax,%rax)
14008f4a0: 4c 89 d9                    	movq	%r11, %rcx
14008f4a3: 31 d2                       	xorl	%edx, %edx
14008f4a5: 4c 89 4c 24 68              	movq	%r9, 0x68(%rsp)
14008f4aa: 4c 89 44 24 60              	movq	%r8, 0x60(%rsp)
14008f4af: e8 9c 61 01 00              	callq	0x1400a5650 <.text+0xa4650>
14008f4b4: 4c 8b 4c 24 68              	movq	0x68(%rsp), %r9
14008f4b9: 4c 8b 44 24 60              	movq	0x60(%rsp), %r8
14008f4be: 48 8b 84 24 18 01 00 00     	movq	0x118(%rsp), %rax
14008f4c6: 48 8b 48 08                 	movq	0x8(%rax), %rcx
14008f4ca: e9 05 ff ff ff              	jmp	0x14008f3d4 <.text+0x8e3d4>
14008f4cf: 90                          	nop
14008f4d0: 48 8b 8c 24 18 01 00 00     	movq	0x118(%rsp), %rcx
14008f4d8: 4c 29 c2                    	subq	%r8, %rdx
14008f4db: e8 20 fb 02 00              	callq	0x1400bf000 <.text+0xbe000>
14008f4e0: e9 b7 fd ff ff              	jmp	0x14008f29c <.text+0x8e29c>
14008f4e5: 4d 85 c9                    	testq	%r9, %r9
14008f4e8: 74 2f                       	je	0x14008f519 <.text+0x8e519>
14008f4ea: 4c 89 d9                    	movq	%r11, %rcx
14008f4ed: 31 d2                       	xorl	%edx, %edx
14008f4ef: 4f 8d 34 0b                 	leaq	(%r11,%r9), %r14
14008f4f3: 4d 89 c8                    	movq	%r9, %r8
14008f4f6: e8 55 61 01 00              	callq	0x1400a5650 <.text+0xa4650>
14008f4fb: 4c 8b 94 24 18 01 00 00     	movq	0x118(%rsp), %r10
14008f503: 8b 54 24 7c                 	movl	0x7c(%rsp), %edx
14008f507: 48 8b 84 24 18 01 00 00     	movq	0x118(%rsp), %rax
14008f50f: 4d 8b 1a                    	movq	(%r10), %r11
14008f512: 48 8b 48 08                 	movq	0x8(%rax), %rcx
14008f516: 4c 89 f0                    	movq	%r14, %rax
14008f519: 48 39 c8                    	cmpq	%rcx, %rax
14008f51c: 74 0c                       	je	0x14008f52a <.text+0x8e52a>
14008f51e: 48 8b 8c 24 18 01 00 00     	movq	0x118(%rsp), %rcx
14008f526: 48 89 41 08                 	movq	%rax, 0x8(%rcx)
14008f52a: 41 89 d6                    	movl	%edx, %r14d
14008f52d: 4c 89 d9                    	movq	%r11, %rcx
14008f530: e9 fe fc ff ff              	jmp	0x14008f233 <.text+0x8e233>
14008f535: 48 89 c6                    	movq	%rax, %rsi
14008f538: e9 ef 91 04 00              	jmp	0x1400d872c <.text+0xd772c>
14008f53d: 48 89 c6                    	movq	%rax, %rsi
14008f540: e9 e7 91 04 00              	jmp	0x1400d872c <.text+0xd772c>
14008f545: 66 66 2e 0f 1f 84 00 00 00 00 00    	nopw	%cs:(%rax,%rax)
14008f550: 41 56                       	pushq	%r14
14008f552: 57                          	pushq	%rdi
14008f553: 56                          	pushq	%rsi
14008f554: 53                          	pushq	%rbx
14008f555: 48 81 ec f8 00 00 00        	subq	$0xf8, %rsp
14008f55c: 48 8b 02                    	movq	(%rdx), %rax
14008f55f: c6 00 00                    	movb	$0x0, (%rax)
14008f562: 48 8b 41 08                 	movq	0x8(%rcx), %rax
14008f566: 48 c7 42 08 00 00 00 00     	movq	$0x0, 0x8(%rdx)
14008f56e: 49 89 d3                    	movq	%rdx, %r11
14008f571: 48 8b 11                    	movq	(%rcx), %rdx
14008f574: 49 89 ca                    	movq	%rcx, %r10
14008f577: 48 29 d0                    	subq	%rdx, %rax
14008f57a: 48 83 f8 46                 	cmpq	$0x46, %rax
14008f57e: 7e 08                       	jle	0x14008f588 <.text+0x8e588>
14008f580: 81 3a d8 3a 91 c7           	cmpl	$0xc7913ad8, (%rdx)     # imm = 0xC7913AD8
14008f586: 74 18                       	je	0x14008f5a0 <.text+0x8e5a0>
14008f588: 45 31 c0                    	xorl	%r8d, %r8d
14008f58b: 44 89 c0                    	movl	%r8d, %eax
14008f58e: 48 81 c4 f8 00 00 00        	addq	$0xf8, %rsp
14008f595: 5b                          	popq	%rbx
14008f596: 5e                          	popq	%rsi
14008f597: 5f                          	popq	%rdi
14008f598: 41 5e                       	popq	%r14
14008f59a: c3                          	retq
14008f59b: 0f 1f 44 00 00              	nopl	(%rax,%rax)
14008f5a0: 80 7a 04 02                 	cmpb	$0x2, 0x4(%rdx)
14008f5a4: 75 e2                       	jne	0x14008f588 <.text+0x8e588>
14008f5a6: 0f b7 7a 15                 	movzwl	0x15(%rdx), %edi
14008f5aa: 48 8d 57 37                 	leaq	0x37(%rdi), %rdx
14008f5ae: 48 39 d0                    	cmpq	%rdx, %rax
14008f5b1: 72 d5                       	jb	0x14008f588 <.text+0x8e588>
14008f5b3: 4c 8d 84 24 b0 00 00 00     	leaq	0xb0(%rsp), %r8
14008f5bb: 48 8d 94 24 90 00 00 00     	leaq	0x90(%rsp), %rdx
14008f5c3: 4c 89 94 24 20 01 00 00     	movq	%r10, 0x120(%rsp)
14008f5cb: 48 8d 4c 24 70              	leaq	0x70(%rsp), %rcx
14008f5d0: 4c 89 9c 24 28 01 00 00     	movq	%r11, 0x128(%rsp)
14008f5d8: 89 7c 24 38                 	movl	%edi, 0x38(%rsp)
14008f5dc: e8 1f ed ff ff              	callq	0x14008e300 <.text+0x8d300>
14008f5e1: 4c 8b 94 24 20 01 00 00     	movq	0x120(%rsp), %r10
14008f5e9: 0f b7 44 24 38              	movzwl	0x38(%rsp), %eax
14008f5ee: 4c 8d 8c 24 d0 00 00 00     	leaq	0xd0(%rsp), %r9
14008f5f6: 48 8d 8c 24 b0 00 00 00     	leaq	0xb0(%rsp), %rcx
14008f5fe: 49 8b 12                    	movq	(%r10), %rdx
14008f601: 44 8d 40 17                 	leal	0x17(%rax), %r8d
14008f605: 89 c6                       	movl	%eax, %esi
14008f607: e8 34 eb ff ff              	callq	0x14008e140 <.text+0x8d140>
14008f60c: 4c 8b 94 24 20 01 00 00     	movq	0x120(%rsp), %r10
14008f614: 84 c0                       	testb	%al, %al
14008f616: 41 89 c0                    	movl	%eax, %r8d
14008f619: 0f 84 91 02 00 00           	je	0x14008f8b0 <.text+0x8e8b0>
14008f61f: 4d 8b 02                    	movq	(%r10), %r8
14008f622: 48 8b 84 24 d0 00 00 00     	movq	0xd0(%rsp), %rax
14008f62a: 48 8b 94 24 d8 00 00 00     	movq	0xd8(%rsp), %rdx
14008f632: 49 8d 4c 38 17              	leaq	0x17(%r8,%rdi), %rcx
14008f637: 48 33 01                    	xorq	(%rcx), %rax
14008f63a: 48 33 51 08                 	xorq	0x8(%rcx), %rdx
14008f63e: 48 09 d0                    	orq	%rdx, %rax
14008f641: 74 1d                       	je	0x14008f660 <.text+0x8e660>
14008f643: b9 20 00 00 00              	movl	$0x20, %ecx
14008f648: 31 c0                       	xorl	%eax, %eax
14008f64a: 48 8d 7c 24 70              	leaq	0x70(%rsp), %rdi
14008f64f: f3 aa                       	rep		stosb	%al, %es:(%rdi)
14008f651: e9 32 ff ff ff              	jmp	0x14008f588 <.text+0x8e588>
14008f656: 66 2e 0f 1f 84 00 00 00 00 00       	nopw	%cs:(%rax,%rax)
14008f660: 48 8b 84 24 e0 00 00 00     	movq	0xe0(%rsp), %rax
14008f668: 48 8b 94 24 e8 00 00 00     	movq	0xe8(%rsp), %rdx
14008f670: 48 33 41 10                 	xorq	0x10(%rcx), %rax
14008f674: 48 33 51 18                 	xorq	0x18(%rcx), %rdx
14008f678: 48 09 d0                    	orq	%rdx, %rax
14008f67b: 75 c6                       	jne	0x14008f643 <.text+0x8e643>
14008f67d: 48 8d 44 24 50              	leaq	0x50(%rsp), %rax
14008f682: f3 41 0f 6f 40 05           	movdqu	0x5(%r8), %xmm0
14008f688: 89 74 24 20                 	movl	%esi, 0x20(%rsp)
14008f68c: 4d 8d 48 17                 	leaq	0x17(%r8), %r9
14008f690: 48 89 44 24 28              	movq	%rax, 0x28(%rsp)
14008f695: 4c 8d 44 24 40              	leaq	0x40(%rsp), %r8
14008f69a: 48 8d 54 24 70              	leaq	0x70(%rsp), %rdx
14008f69f: 31 c9                       	xorl	%ecx, %ecx
14008f6a1: 0f 29 44 24 40              	movaps	%xmm0, 0x40(%rsp)
14008f6a6: 66 0f ef c0                 	pxor	%xmm0, %xmm0
14008f6aa: 48 c7 44 24 60 00 00 00 00  	movq	$0x0, 0x60(%rsp)
14008f6b3: 0f 29 44 24 50              	movaps	%xmm0, 0x50(%rsp)
14008f6b8: e8 33 f9 ff ff              	callq	0x14008eff0 <.text+0x8dff0>
14008f6bd: 41 89 c0                    	movl	%eax, %r8d
14008f6c0: 48 8b 44 24 50              	movq	0x50(%rsp), %rax
14008f6c5: 4c 8b 9c 24 28 01 00 00     	movq	0x128(%rsp), %r11
14008f6cd: 4c 8b 74 24 60              	movq	0x60(%rsp), %r14
14008f6d2: 48 89 c3                    	movq	%rax, %rbx
14008f6d5: 45 84 c0                    	testb	%r8b, %r8b
14008f6d8: 0f 84 79 01 00 00           	je	0x14008f857 <.text+0x8e857>
14008f6de: 48 8b 54 24 58              	movq	0x58(%rsp), %rdx
14008f6e3: 48 89 d7                    	movq	%rdx, %rdi
14008f6e6: 48 29 c7                    	subq	%rax, %rdi
14008f6e9: 48 89 fe                    	movq	%rdi, %rsi
14008f6ec: 74 4d                       	je	0x14008f73b <.text+0x8e73b>
14008f6ee: b9 03 00 00 00              	movl	$0x3, %ecx
14008f6f3: 45 31 c9                    	xorl	%r9d, %r9d
14008f6f6: 31 c0                       	xorl	%eax, %eax
14008f6f8: 0f 1f 84 00 00 00 00 00     	nopl	(%rax,%rax)
14008f700: 49 89 c0                    	movq	%rax, %r8
14008f703: 49 89 ca                    	movq	%rcx, %r10
14008f706: 48 83 c1 07                 	addq	$0x7, %rcx
14008f70a: 41 83 e0 1f                 	andl	$0x1f, %r8d
14008f70e: 41 83 e2 1f                 	andl	$0x1f, %r10d
14008f712: 42 0f b6 bc 04 90 00 00 00  	movzbl	0x90(%rsp,%r8), %edi
14008f71b: 46 0f b6 84 14 90 00 00 00  	movzbl	0x90(%rsp,%r10), %r8d
14008f724: 44 01 cf                    	addl	%r9d, %edi
14008f727: 41 83 c1 11                 	addl	$0x11, %r9d
14008f72b: 41 31 f8                    	xorl	%edi, %r8d
14008f72e: 44 30 04 03                 	xorb	%r8b, (%rbx,%rax)
14008f732: 48 83 c0 01                 	addq	$0x1, %rax
14008f736: 48 39 c6                    	cmpq	%rax, %rsi
14008f739: 75 c5                       	jne	0x14008f700 <.text+0x8e700>
14008f73b: 48 39 d3                    	cmpq	%rdx, %rbx
14008f73e: 0f 84 c1 00 00 00           	je	0x14008f805 <.text+0x8e805>
14008f744: 48 8d 4a ff                 	leaq	-0x1(%rdx), %rcx
14008f748: 48 89 d8                    	movq	%rbx, %rax
14008f74b: 48 29 d9                    	subq	%rbx, %rcx
14008f74e: 48 83 f9 0e                 	cmpq	$0xe, %rcx
14008f752: 76 63                       	jbe	0x14008f7b7 <.text+0x8e7b7>
14008f754: 49 89 d0                    	movq	%rdx, %r8
14008f757: bf 0f 0f 0f 0f              	movl	$0xf0f0f0f, %edi        # imm = 0xF0F0F0F
14008f75c: 49 29 d8                    	subq	%rbx, %r8
14008f75f: 66 0f 6e d7                 	movd	%edi, %xmm2
14008f763: 4d 89 c1                    	movq	%r8, %r9
14008f766: 66 0f 70 d2 00              	pshufd	$0x0, %xmm2, %xmm2      # xmm2 = xmm2[0,0,0,0]
14008f76b: 49 83 e1 f0                 	andq	$-0x10, %r9
14008f76f: 4a 8d 0c 0b                 	leaq	(%rbx,%r9), %rcx
14008f773: 66 90                       	nop
14008f775: 66 66 2e 0f 1f 84 00 00 00 00 00    	nopw	%cs:(%rax,%rax)
14008f780: f3 0f 6f 00                 	movdqu	(%rax), %xmm0
14008f784: 48 83 c0 10                 	addq	$0x10, %rax
14008f788: 66 0f 6f c8                 	movdqa	%xmm0, %xmm1
14008f78c: 66 0f fc c0                 	paddb	%xmm0, %xmm0
14008f790: 66 0f fc c0                 	paddb	%xmm0, %xmm0
14008f794: 66 0f 71 d1 04              	psrlw	$0x4, %xmm1
14008f799: 66 0f fc c0                 	paddb	%xmm0, %xmm0
14008f79d: 66 0f db ca                 	pand	%xmm2, %xmm1
14008f7a1: 66 0f fc c0                 	paddb	%xmm0, %xmm0
14008f7a5: 66 0f eb c1                 	por	%xmm1, %xmm0
14008f7a9: 0f 11 40 f0                 	movups	%xmm0, -0x10(%rax)
14008f7ad: 48 39 c8                    	cmpq	%rcx, %rax
14008f7b0: 75 ce                       	jne	0x14008f780 <.text+0x8e780>
14008f7b2: 4d 39 c1                    	cmpq	%r8, %r9
14008f7b5: 74 42                       	je	0x14008f7f9 <.text+0x8e7f9>
14008f7b7: 48 89 d1                    	movq	%rdx, %rcx
14008f7ba: 48 29 c1                    	subq	%rax, %rcx
14008f7bd: 83 e1 01                    	andl	$0x1, %ecx
14008f7c0: 74 0e                       	je	0x14008f7d0 <.text+0x8e7d0>
14008f7c2: c0 00 04                    	rolb	$0x4, (%rax)
14008f7c5: 48 83 c0 01                 	addq	$0x1, %rax
14008f7c9: 48 39 d0                    	cmpq	%rdx, %rax
14008f7cc: 74 2b                       	je	0x14008f7f9 <.text+0x8e7f9>
14008f7ce: 66 90                       	nop
14008f7d0: c0 00 04                    	rolb	$0x4, (%rax)
14008f7d3: c0 40 01 04                 	rolb	$0x4, 0x1(%rax)
14008f7d7: 48 83 c0 02                 	addq	$0x2, %rax
14008f7db: 48 39 d0                    	cmpq	%rdx, %rax
14008f7de: 75 f0                       	jne	0x14008f7d0 <.text+0x8e7d0>
14008f7e0: eb 17                       	jmp	0x14008f7f9 <.text+0x8e7f9>
14008f7e2: 0f 1f 00                    	nopl	(%rax)
14008f7e5: 66 66 2e 0f 1f 84 00 00 00 00 00    	nopw	%cs:(%rax,%rax)
14008f7f0: 48 83 ea 01                 	subq	$0x1, %rdx
14008f7f4: 48 39 da                    	cmpq	%rbx, %rdx
14008f7f7: 74 06                       	je	0x14008f7ff <.text+0x8e7ff>
14008f7f9: 80 7a ff 00                 	cmpb	$0x0, -0x1(%rdx)
14008f7fd: 74 f1                       	je	0x14008f7f0 <.text+0x8e7f0>
14008f7ff: 48 29 da                    	subq	%rbx, %rdx
14008f802: 48 89 d6                    	movq	%rdx, %rsi
14008f805: 48 89 74 24 20              	movq	%rsi, 0x20(%rsp)
14008f80a: 4d 8b 43 08                 	movq	0x8(%r11), %r8
14008f80e: 49 89 d9                    	movq	%rbx, %r9
14008f811: 31 d2                       	xorl	%edx, %edx
14008f813: 4c 89 d9                    	movq	%r11, %rcx
14008f816: 4c 89 9c 24 28 01 00 00     	movq	%r11, 0x128(%rsp)
14008f81e: e8 ad 00 03 00              	callq	0x1400bf8d0 <.text+0xbe8d0>
14008f823: 48 8b 8c 24 28 01 00 00     	movq	0x128(%rsp), %rcx
14008f82b: e8 a0 f1 ff ff              	callq	0x14008e9d0 <.text+0x8d9d0>
14008f830: 41 89 c0                    	movl	%eax, %r8d
14008f833: 84 c0                       	testb	%al, %al
14008f835: 75 16                       	jne	0x14008f84d <.text+0x8e84d>
14008f837: 4c 8b 9c 24 28 01 00 00     	movq	0x128(%rsp), %r11
14008f83f: 49 8b 03                    	movq	(%r11), %rax
14008f842: c6 00 00                    	movb	$0x0, (%rax)
14008f845: 49 c7 43 08 00 00 00 00     	movq	$0x0, 0x8(%r11)
14008f84d: 48 89 f1                    	movq	%rsi, %rcx
14008f850: 31 c0                       	xorl	%eax, %eax
14008f852: 48 89 df                    	movq	%rbx, %rdi
14008f855: f3 aa                       	rep		stosb	%al, %es:(%rdi)
14008f857: 31 c0                       	xorl	%eax, %eax
14008f859: b9 20 00 00 00              	movl	$0x20, %ecx
14008f85e: 48 8d 7c 24 70              	leaq	0x70(%rsp), %rdi
14008f863: f3 aa                       	rep		stosb	%al, %es:(%rdi)
14008f865: b9 20 00 00 00              	movl	$0x20, %ecx
14008f86a: 48 8d bc 24 90 00 00 00     	leaq	0x90(%rsp), %rdi
14008f872: f3 aa                       	rep		stosb	%al, %es:(%rdi)
14008f874: b9 20 00 00 00              	movl	$0x20, %ecx
14008f879: 48 8d bc 24 b0 00 00 00     	leaq	0xb0(%rsp), %rdi
14008f881: f3 aa                       	rep		stosb	%al, %es:(%rdi)
14008f883: 48 85 db                    	testq	%rbx, %rbx
14008f886: 0f 84 ff fc ff ff           	je	0x14008f58b <.text+0x8e58b>
14008f88c: 4c 89 f2                    	movq	%r14, %rdx
14008f88f: 48 89 d9                    	movq	%rbx, %rcx
14008f892: 44 88 44 24 38              	movb	%r8b, 0x38(%rsp)
14008f897: 48 29 da                    	subq	%rbx, %rdx
14008f89a: e8 31 a8 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008f89f: 44 0f b6 44 24 38           	movzbl	0x38(%rsp), %r8d
14008f8a5: e9 e1 fc ff ff              	jmp	0x14008f58b <.text+0x8e58b>
14008f8aa: 66 0f 1f 44 00 00           	nopw	(%rax,%rax)
14008f8b0: b9 20 00 00 00              	movl	$0x20, %ecx
14008f8b5: 48 8d 7c 24 70              	leaq	0x70(%rsp), %rdi
14008f8ba: f3 aa                       	rep		stosb	%al, %es:(%rdi)
14008f8bc: e9 ca fc ff ff              	jmp	0x14008f58b <.text+0x8e58b>
14008f8c1: e9 86 8e 04 00              	jmp	0x1400d874c <.text+0xd774c>
14008f8c6: e9 77 8e 04 00              	jmp	0x1400d8742 <.text+0xd7742>
14008f8cb: 0f 1f 44 00 00              	nopl	(%rax,%rax)
14008f8d0: 41 57                       	pushq	%r15
14008f8d2: 41 56                       	pushq	%r14
14008f8d4: 41 55                       	pushq	%r13
14008f8d6: 41 54                       	pushq	%r12
14008f8d8: 55                          	pushq	%rbp
14008f8d9: 57                          	pushq	%rdi
14008f8da: 56                          	pushq	%rsi
14008f8db: 53                          	pushq	%rbx
14008f8dc: 48 81 ec 68 01 00 00        	subq	$0x168, %rsp            # imm = 0x168
14008f8e3: 48 8d 15 12 be 06 00        	leaq	0x6be12(%rip), %rdx     # 0x1400fb6fc
14008f8ea: 4c 8d 42 18                 	leaq	0x18(%rdx), %r8
14008f8ee: 4c 8d ac 24 d0 00 00 00     	leaq	0xd0(%rsp), %r13
14008f8f6: 48 8d 8c 24 c0 00 00 00     	leaq	0xc0(%rsp), %rcx
14008f8fe: 4c 89 ac 24 c0 00 00 00     	movq	%r13, 0xc0(%rsp)
14008f906: e8 95 eb ff ff              	callq	0x14008e4a0 <.text+0x8d4a0>
14008f90b: 48 8d 8c 24 c0 00 00 00     	leaq	0xc0(%rsp), %rcx
14008f913: e8 b8 f0 ff ff              	callq	0x14008e9d0 <.text+0x8d9d0>
14008f918: 89 c5                       	movl	%eax, %ebp
14008f91a: 84 c0                       	testb	%al, %al
14008f91c: 75 3a                       	jne	0x14008f958 <.text+0x8e958>
14008f91e: 48 8b 8c 24 c0 00 00 00     	movq	0xc0(%rsp), %rcx
14008f926: 4c 39 e9                    	cmpq	%r13, %rcx
14008f929: 74 11                       	je	0x14008f93c <.text+0x8e93c>
14008f92b: 48 8b 84 24 d0 00 00 00     	movq	0xd0(%rsp), %rax
14008f933: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008f937: e8 94 a7 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008f93c: 89 e8                       	movl	%ebp, %eax
14008f93e: 48 81 c4 68 01 00 00        	addq	$0x168, %rsp            # imm = 0x168
14008f945: 5b                          	popq	%rbx
14008f946: 5e                          	popq	%rsi
14008f947: 5f                          	popq	%rdi
14008f948: 5d                          	popq	%rbp
14008f949: 41 5c                       	popq	%r12
14008f94b: 41 5d                       	popq	%r13
14008f94d: 41 5e                       	popq	%r14
14008f94f: 41 5f                       	popq	%r15
14008f951: c3                          	retq
14008f952: 66 0f 1f 44 00 00           	nopw	(%rax,%rax)
14008f958: 4c 8d 84 24 20 01 00 00     	leaq	0x120(%rsp), %r8
14008f960: 48 8d 94 24 00 01 00 00     	leaq	0x100(%rsp), %rdx
14008f968: 48 8d 8c 24 e0 00 00 00     	leaq	0xe0(%rsp), %rcx
14008f970: e8 8b e9 ff ff              	callq	0x14008e300 <.text+0x8d300>
14008f975: 4c 8b a4 24 c8 00 00 00     	movq	0xc8(%rsp), %r12
14008f97d: 4d 85 e4                    	testq	%r12, %r12
14008f980: 0f 84 7a 04 00 00           	je	0x14008fe00 <.text+0x8ee00>
14008f986: 4c 89 e1                    	movq	%r12, %rcx
14008f989: 48 8b 9c 24 c0 00 00 00     	movq	0xc0(%rsp), %rbx
14008f991: e8 4a a7 03 00              	callq	0x1400ca0e0 <.text+0xc90e0>
14008f996: 4d 89 e0                    	movq	%r12, %r8
14008f999: 48 89 da                    	movq	%rbx, %rdx
14008f99c: 48 89 c1                    	movq	%rax, %rcx
14008f99f: 48 89 c6                    	movq	%rax, %rsi
14008f9a2: 4a 8d 3c 20                 	leaq	(%rax,%r12), %rdi
14008f9a6: e8 3d 5e 01 00              	callq	0x1400a57e8 <.text+0xa47e8>
14008f9ab: 49 8d 44 24 ff              	leaq	-0x1(%r12), %rax
14008f9b0: 48 83 f8 0e                 	cmpq	$0xe, %rax
14008f9b4: 48 89 f0                    	movq	%rsi, %rax
14008f9b7: 76 7e                       	jbe	0x14008fa37 <.text+0x8ea37>
14008f9b9: 4c 89 e1                    	movq	%r12, %rcx
14008f9bc: bb 0f 0f 0f 0f              	movl	$0xf0f0f0f, %ebx        # imm = 0xF0F0F0F
14008f9c1: 48 83 e1 f0                 	andq	$-0x10, %rcx
14008f9c5: 66 0f 6e d3                 	movd	%ebx, %xmm2
14008f9c9: 48 8d 14 0e                 	leaq	(%rsi,%rcx), %rdx
14008f9cd: 66 0f 70 d2 00              	pshufd	$0x0, %xmm2, %xmm2      # xmm2 = xmm2[0,0,0,0]
14008f9d2: 66 90                       	nop
14008f9d4: 66 66 2e 0f 1f 84 00 00 00 00 00    	nopw	%cs:(%rax,%rax)
14008f9df: 66 66 2e 0f 1f 84 00 00 00 00 00    	nopw	%cs:(%rax,%rax)
14008f9ea: 66 66 2e 0f 1f 84 00 00 00 00 00    	nopw	%cs:(%rax,%rax)
14008f9f5: 66 66 2e 0f 1f 84 00 00 00 00 00    	nopw	%cs:(%rax,%rax)
14008fa00: f3 0f 6f 08                 	movdqu	(%rax), %xmm1
14008fa04: 48 83 c0 10                 	addq	$0x10, %rax
14008fa08: 66 0f 6f c1                 	movdqa	%xmm1, %xmm0
14008fa0c: 66 0f fc c1                 	paddb	%xmm1, %xmm0
14008fa10: 66 0f 71 d1 04              	psrlw	$0x4, %xmm1
14008fa15: 66 0f fc c0                 	paddb	%xmm0, %xmm0
14008fa19: 66 0f db ca                 	pand	%xmm2, %xmm1
14008fa1d: 66 0f fc c0                 	paddb	%xmm0, %xmm0
14008fa21: 66 0f fc c0                 	paddb	%xmm0, %xmm0
14008fa25: 66 0f eb c1                 	por	%xmm1, %xmm0
14008fa29: 0f 11 40 f0                 	movups	%xmm0, -0x10(%rax)
14008fa2d: 48 39 d0                    	cmpq	%rdx, %rax
14008fa30: 75 ce                       	jne	0x14008fa00 <.text+0x8ea00>
14008fa32: 49 39 cc                    	cmpq	%rcx, %r12
14008fa35: 74 29                       	je	0x14008fa60 <.text+0x8ea60>
14008fa37: 48 89 fa                    	movq	%rdi, %rdx
14008fa3a: 48 29 c2                    	subq	%rax, %rdx
14008fa3d: 83 e2 01                    	andl	$0x1, %edx
14008fa40: 74 0e                       	je	0x14008fa50 <.text+0x8ea50>
14008fa42: c0 00 04                    	rolb	$0x4, (%rax)
14008fa45: 48 83 c0 01                 	addq	$0x1, %rax
14008fa49: 48 39 f8                    	cmpq	%rdi, %rax
14008fa4c: 74 12                       	je	0x14008fa60 <.text+0x8ea60>
14008fa4e: 66 90                       	nop
14008fa50: c0 00 04                    	rolb	$0x4, (%rax)
14008fa53: c0 40 01 04                 	rolb	$0x4, 0x1(%rax)
14008fa57: 48 83 c0 02                 	addq	$0x2, %rax
14008fa5b: 48 39 f8                    	cmpq	%rdi, %rax
14008fa5e: 75 f0                       	jne	0x14008fa50 <.text+0x8ea50>
14008fa60: 4c 89 e5                    	movq	%r12, %rbp
14008fa63: b9 03 00 00 00              	movl	$0x3, %ecx
14008fa68: 31 d2                       	xorl	%edx, %edx
14008fa6a: 31 db                       	xorl	%ebx, %ebx
14008fa6c: 66 0f 1f 84 00 00 00 00 00  	nopw	(%rax,%rax)
14008fa75: 66 66 2e 0f 1f 84 00 00 00 00 00    	nopw	%cs:(%rax,%rax)
14008fa80: 48 89 d8                    	movq	%rbx, %rax
14008fa83: 49 89 c8                    	movq	%rcx, %r8
14008fa86: 48 83 c1 07                 	addq	$0x7, %rcx
14008fa8a: 83 e0 1f                    	andl	$0x1f, %eax
14008fa8d: 41 83 e0 1f                 	andl	$0x1f, %r8d
14008fa91: 44 0f b6 94 04 00 01 00 00  	movzbl	0x100(%rsp,%rax), %r10d
14008fa9a: 42 0f b6 84 04 00 01 00 00  	movzbl	0x100(%rsp,%r8), %eax
14008faa3: 41 01 d2                    	addl	%edx, %r10d
14008faa6: 83 c2 11                    	addl	$0x11, %edx
14008faa9: 44 31 d0                    	xorl	%r10d, %eax
14008faac: 30 04 1e                    	xorb	%al, (%rsi,%rbx)
14008faaf: 48 83 c3 01                 	addq	$0x1, %rbx
14008fab3: 49 39 dc                    	cmpq	%rbx, %r12
14008fab6: 75 c8                       	jne	0x14008fa80 <.text+0x8ea80>
14008fab8: 41 b9 02 00 00 00           	movl	$0x2, %r9d
14008fabe: 41 b8 10 00 00 00           	movl	$0x10, %r8d
14008fac4: 48 8d 54 24 60              	leaq	0x60(%rsp), %rdx
14008fac9: 31 c9                       	xorl	%ecx, %ecx
14008facb: e8 c0 11 00 00              	callq	0x140090c90 <.text+0x8fc90>
14008fad0: 85 c0                       	testl	%eax, %eax
14008fad2: 0f 88 c0 01 00 00           	js	0x14008fc98 <.text+0x8ec98>
14008fad8: 89 6c 24 20                 	movl	%ebp, 0x20(%rsp)
14008fadc: 66 0f ef c0                 	pxor	%xmm0, %xmm0
14008fae0: 49 89 f1                    	movq	%rsi, %r9
14008fae3: 4c 8d 44 24 60              	leaq	0x60(%rsp), %r8
14008fae8: 48 8d 84 24 a0 00 00 00     	leaq	0xa0(%rsp), %rax
14008faf0: 48 8d 94 24 e0 00 00 00     	leaq	0xe0(%rsp), %rdx
14008faf8: b9 01 00 00 00              	movl	$0x1, %ecx
14008fafd: 48 c7 84 24 b0 00 00 00 00 00 00 00 	movq	$0x0, 0xb0(%rsp)
14008fb09: 48 89 44 24 28              	movq	%rax, 0x28(%rsp)
14008fb0e: 0f 29 84 24 a0 00 00 00     	movaps	%xmm0, 0xa0(%rsp)
14008fb16: e8 d5 f4 ff ff              	callq	0x14008eff0 <.text+0x8dff0>
14008fb1b: 48 8b 94 24 a0 00 00 00     	movq	0xa0(%rsp), %rdx
14008fb23: 48 8b 8c 24 b0 00 00 00     	movq	0xb0(%rsp), %rcx
14008fb2b: 48 89 54 24 38              	movq	%rdx, 0x38(%rsp)
14008fb30: 48 89 4c 24 48              	movq	%rcx, 0x48(%rsp)
14008fb35: 84 c0                       	testb	%al, %al
14008fb37: 0f 84 83 02 00 00           	je	0x14008fdc0 <.text+0x8edc0>
14008fb3d: 48 8b ac 24 a8 00 00 00     	movq	0xa8(%rsp), %rbp
14008fb45: 48 89 6c 24 50              	movq	%rbp, 0x50(%rsp)
14008fb4a: 48 29 d5                    	subq	%rdx, %rbp
14008fb4d: 48 81 fd ff 7f 00 00        	cmpq	$0x7fff, %rbp           # imm = 0x7FFF
14008fb54: 0f 8f 66 02 00 00           	jg	0x14008fdc0 <.text+0x8edc0>
14008fb5a: b9 04 00 00 00              	movl	$0x4, %ecx
14008fb5f: e8 7c a5 03 00              	callq	0x1400ca0e0 <.text+0xc90e0>
14008fb64: c7 00 d8 3a 91 c7           	movl	$0xc7913ad8, (%rax)     # imm = 0xC7913AD8
14008fb6a: b9 08 00 00 00              	movl	$0x8, %ecx
14008fb6f: 49 89 c4                    	movq	%rax, %r12
14008fb72: e8 69 a5 03 00              	callq	0x1400ca0e0 <.text+0xc90e0>
14008fb77: c6 40 04 02                 	movb	$0x2, 0x4(%rax)
14008fb7b: 49 89 c6                    	movq	%rax, %r14
14008fb7e: 4c 8d 78 05                 	leaq	0x5(%rax), %r15
14008fb82: ba 04 00 00 00              	movl	$0x4, %edx
14008fb87: c7 00 d8 3a 91 c7           	movl	$0xc7913ad8, (%rax)     # imm = 0xC7913AD8
14008fb8d: 4c 89 e1                    	movq	%r12, %rcx
14008fb90: e8 3b a5 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008fb95: 4c 89 b4 24 80 00 00 00     	movq	%r14, 0x80(%rsp)
14008fb9d: 4c 8d 4c 24 70              	leaq	0x70(%rsp), %r9
14008fba2: 4c 89 fa                    	movq	%r15, %rdx
14008fba5: 49 83 c6 08                 	addq	$0x8, %r14
14008fba9: 4c 8d 44 24 60              	leaq	0x60(%rsp), %r8
14008fbae: 48 8d 8c 24 80 00 00 00     	leaq	0x80(%rsp), %rcx
14008fbb6: 4c 89 bc 24 88 00 00 00     	movq	%r15, 0x88(%rsp)
14008fbbe: 4c 89 b4 24 90 00 00 00     	movq	%r14, 0x90(%rsp)
14008fbc6: e8 65 eb ff ff              	callq	0x14008e730 <.text+0x8d730>
14008fbcb: 48 8d 84 24 40 01 00 00     	leaq	0x140(%rsp), %rax
14008fbd3: 48 8d 8c 24 80 00 00 00     	leaq	0x80(%rsp), %rcx
14008fbdb: 40 88 ac 24 40 01 00 00     	movb	%bpl, 0x140(%rsp)
14008fbe3: 48 89 c2                    	movq	%rax, %rdx
14008fbe6: 48 89 44 24 40              	movq	%rax, 0x40(%rsp)
14008fbeb: e8 20 ea ff ff              	callq	0x14008e610 <.text+0x8d610>
14008fbf0: 48 8b 54 24 40              	movq	0x40(%rsp), %rdx
14008fbf5: 48 89 e8                    	movq	%rbp, %rax
14008fbf8: 48 8d 8c 24 80 00 00 00     	leaq	0x80(%rsp), %rcx
14008fc00: 88 a4 24 40 01 00 00        	movb	%ah, 0x140(%rsp)
14008fc07: e8 04 ea ff ff              	callq	0x14008e610 <.text+0x8d610>
14008fc0c: 4c 8b bc 24 88 00 00 00     	movq	0x88(%rsp), %r15
14008fc14: 4c 8b a4 24 80 00 00 00     	movq	0x80(%rsp), %r12
14008fc1c: 4c 8b b4 24 90 00 00 00     	movq	0x90(%rsp), %r14
14008fc24: 48 8b 54 24 38              	movq	0x38(%rsp), %rdx
14008fc29: 4d 89 f8                    	movq	%r15, %r8
14008fc2c: 4d 29 e0                    	subq	%r12, %r8
14008fc2f: 48 39 54 24 50              	cmpq	%rdx, 0x50(%rsp)
14008fc34: 74 28                       	je	0x14008fc5e <.text+0x8ec5e>
14008fc36: 4c 89 f0                    	movq	%r14, %rax
14008fc39: 4c 29 f8                    	subq	%r15, %rax
14008fc3c: 48 39 e8                    	cmpq	%rbp, %rax
14008fc3f: 0f 82 85 03 00 00           	jb	0x14008ffca <.text+0x8efca>
14008fc45: 48 85 ed                    	testq	%rbp, %rbp
14008fc48: 74 0b                       	je	0x14008fc55 <.text+0x8ec55>
14008fc4a: 49 89 e8                    	movq	%rbp, %r8
14008fc4d: 4c 89 f9                    	movq	%r15, %rcx
14008fc50: e8 93 5b 01 00              	callq	0x1400a57e8 <.text+0xa47e8>
14008fc55: 49 01 ef                    	addq	%rbp, %r15
14008fc58: 4d 89 f8                    	movq	%r15, %r8
14008fc5b: 4d 29 e0                    	subq	%r12, %r8
14008fc5e: 4c 8b 4c 24 40              	movq	0x40(%rsp), %r9
14008fc63: 4c 89 e2                    	movq	%r12, %rdx
14008fc66: 48 8d 8c 24 20 01 00 00     	leaq	0x120(%rsp), %rcx
14008fc6e: e8 cd e4 ff ff              	callq	0x14008e140 <.text+0x8d140>
14008fc73: 89 c5                       	movl	%eax, %ebp
14008fc75: 84 c0                       	testb	%al, %al
14008fc77: 0f 85 93 01 00 00           	jne	0x14008fe10 <.text+0x8ee10>
14008fc7d: 48 8d bc 24 e0 00 00 00     	leaq	0xe0(%rsp), %rdi
14008fc85: b9 20 00 00 00              	movl	$0x20, %ecx
14008fc8a: f3 aa                       	rep		stosb	%al, %es:(%rdi)
14008fc8c: e9 4b 01 00 00              	jmp	0x14008fddc <.text+0x8eddc>
14008fc91: 0f 1f 80 00 00 00 00        	nopl	(%rax)
14008fc98: 48 8d bc 24 e0 00 00 00     	leaq	0xe0(%rsp), %rdi
14008fca0: b9 20 00 00 00              	movl	$0x20, %ecx
14008fca5: 31 c0                       	xorl	%eax, %eax
14008fca7: f3 aa                       	rep		stosb	%al, %es:(%rdi)
14008fca9: 45 31 f6                    	xorl	%r14d, %r14d
14008fcac: 45 31 ff                    	xorl	%r15d, %r15d
14008fcaf: 45 31 e4                    	xorl	%r12d, %r12d
14008fcb2: 31 ed                       	xorl	%ebp, %ebp
14008fcb4: 48 85 f6                    	testq	%rsi, %rsi
14008fcb7: 74 0b                       	je	0x14008fcc4 <.text+0x8ecc4>
14008fcb9: 48 89 da                    	movq	%rbx, %rdx
14008fcbc: 48 89 f1                    	movq	%rsi, %rcx
14008fcbf: e8 0c a4 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008fcc4: 48 8b 8c 24 c0 00 00 00     	movq	0xc0(%rsp), %rcx
14008fccc: 4c 39 e9                    	cmpq	%r13, %rcx
14008fccf: 74 11                       	je	0x14008fce2 <.text+0x8ece2>
14008fcd1: 48 8b 84 24 d0 00 00 00     	movq	0xd0(%rsp), %rax
14008fcd9: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008fcdd: e8 ee a3 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008fce2: 40 84 ed                    	testb	%bpl, %bpl
14008fce5: 0f 84 82 00 00 00           	je	0x14008fd6d <.text+0x8ed6d>
14008fceb: 48 8d 9c 24 10 01 00 00     	leaq	0x110(%rsp), %rbx
14008fcf3: 48 8d 8c 24 80 00 00 00     	leaq	0x80(%rsp), %rcx
14008fcfb: c6 84 24 10 01 00 00 00     	movb	$0x0, 0x110(%rsp)
14008fd03: 48 8d 94 24 00 01 00 00     	leaq	0x100(%rsp), %rdx
14008fd0b: 48 89 9c 24 00 01 00 00     	movq	%rbx, 0x100(%rsp)
14008fd13: 48 c7 84 24 08 01 00 00 00 00 00 00 	movq	$0x0, 0x108(%rsp)
14008fd1f: 4c 89 a4 24 80 00 00 00     	movq	%r12, 0x80(%rsp)
14008fd27: 4c 89 bc 24 88 00 00 00     	movq	%r15, 0x88(%rsp)
14008fd2f: 4c 89 b4 24 90 00 00 00     	movq	%r14, 0x90(%rsp)
14008fd37: e8 14 f8 ff ff              	callq	0x14008f550 <.text+0x8e550>
14008fd3c: 89 c5                       	movl	%eax, %ebp
14008fd3e: 84 c0                       	testb	%al, %al
14008fd40: 74 0d                       	je	0x14008fd4f <.text+0x8ed4f>
14008fd42: 48 83 bc 24 08 01 00 00 18  	cmpq	$0x18, 0x108(%rsp)
14008fd4b: 74 43                       	je	0x14008fd90 <.text+0x8ed90>
14008fd4d: 31 ed                       	xorl	%ebp, %ebp
14008fd4f: 48 8b 8c 24 00 01 00 00     	movq	0x100(%rsp), %rcx
14008fd57: 48 39 d9                    	cmpq	%rbx, %rcx
14008fd5a: 74 11                       	je	0x14008fd6d <.text+0x8ed6d>
14008fd5c: 48 8b 84 24 10 01 00 00     	movq	0x110(%rsp), %rax
14008fd64: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008fd68: e8 63 a3 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008fd6d: 4d 85 e4                    	testq	%r12, %r12
14008fd70: 0f 84 c6 fb ff ff           	je	0x14008f93c <.text+0x8e93c>
14008fd76: 4c 89 f2                    	movq	%r14, %rdx
14008fd79: 4c 89 e1                    	movq	%r12, %rcx
14008fd7c: 4c 29 e2                    	subq	%r12, %rdx
14008fd7f: e8 4c a3 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008fd84: e9 b3 fb ff ff              	jmp	0x14008f93c <.text+0x8e93c>
14008fd89: 0f 1f 80 00 00 00 00        	nopl	(%rax)
14008fd90: 48 8b 8c 24 00 01 00 00     	movq	0x100(%rsp), %rcx
14008fd98: 48 b8 68 74 74 70 3a 2f 2f 32       	movabsq	$0x322f2f3a70747468, %rax # imm = 0x322F2F3A70747468
14008fda2: 48 ba 30 33 2e 30 2e 31 31 33       	movabsq	$0x3331312e302e3330, %rdx # imm = 0x3331312E302E3330
14008fdac: 48 33 01                    	xorq	(%rcx), %rax
14008fdaf: 48 33 51 08                 	xorq	0x8(%rcx), %rdx
14008fdb3: 48 09 d0                    	orq	%rdx, %rax
14008fdb6: 0f 84 25 01 00 00           	je	0x14008fee1 <.text+0x8eee1>
14008fdbc: 31 ed                       	xorl	%ebp, %ebp
14008fdbe: eb 9c                       	jmp	0x14008fd5c <.text+0x8ed5c>
14008fdc0: 48 8d bc 24 e0 00 00 00     	leaq	0xe0(%rsp), %rdi
14008fdc8: b9 20 00 00 00              	movl	$0x20, %ecx
14008fdcd: 31 c0                       	xorl	%eax, %eax
14008fdcf: f3 aa                       	rep		stosb	%al, %es:(%rdi)
14008fdd1: 45 31 f6                    	xorl	%r14d, %r14d
14008fdd4: 45 31 ff                    	xorl	%r15d, %r15d
14008fdd7: 45 31 e4                    	xorl	%r12d, %r12d
14008fdda: 31 ed                       	xorl	%ebp, %ebp
14008fddc: 48 8b 4c 24 38              	movq	0x38(%rsp), %rcx
14008fde1: 48 85 c9                    	testq	%rcx, %rcx
14008fde4: 0f 84 ca fe ff ff           	je	0x14008fcb4 <.text+0x8ecb4>
14008fdea: 48 8b 54 24 48              	movq	0x48(%rsp), %rdx
14008fdef: 48 29 ca                    	subq	%rcx, %rdx
14008fdf2: e8 d9 a2 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008fdf7: e9 b8 fe ff ff              	jmp	0x14008fcb4 <.text+0x8ecb4>
14008fdfc: 0f 1f 40 00                 	nopl	(%rax)
14008fe00: 31 db                       	xorl	%ebx, %ebx
14008fe02: 31 ff                       	xorl	%edi, %edi
14008fe04: 31 f6                       	xorl	%esi, %esi
14008fe06: 31 ed                       	xorl	%ebp, %ebp
14008fe08: e9 ab fc ff ff              	jmp	0x14008fab8 <.text+0x8eab8>
14008fe0d: 0f 1f 00                    	nopl	(%rax)
14008fe10: 4c 8b 44 24 40              	movq	0x40(%rsp), %r8
14008fe15: 4c 8d 8c 24 60 01 00 00     	leaq	0x160(%rsp), %r9
14008fe1d: 4c 89 fa                    	movq	%r15, %rdx
14008fe20: 48 8d 8c 24 80 00 00 00     	leaq	0x80(%rsp), %rcx
14008fe28: 4c 89 a4 24 80 00 00 00     	movq	%r12, 0x80(%rsp)
14008fe30: 4c 89 bc 24 88 00 00 00     	movq	%r15, 0x88(%rsp)
14008fe38: 4c 89 b4 24 90 00 00 00     	movq	%r14, 0x90(%rsp)
14008fe40: e8 eb e8 ff ff              	callq	0x14008e730 <.text+0x8d730>
14008fe45: 41 b9 02 00 00 00           	movl	$0x2, %r9d
14008fe4b: 41 b8 10 00 00 00           	movl	$0x10, %r8d
14008fe51: 48 8d 54 24 70              	leaq	0x70(%rsp), %rdx
14008fe56: 31 c9                       	xorl	%ecx, %ecx
14008fe58: 4c 8b a4 24 80 00 00 00     	movq	0x80(%rsp), %r12
14008fe60: 4c 8b bc 24 88 00 00 00     	movq	0x88(%rsp), %r15
14008fe68: 4c 8b b4 24 90 00 00 00     	movq	0x90(%rsp), %r14
14008fe70: e8 1b 0e 00 00              	callq	0x140090c90 <.text+0x8fc90>
14008fe75: 4c 8d 8c 24 80 00 00 00     	leaq	0x80(%rsp), %r9
14008fe7d: 4c 8d 44 24 70              	leaq	0x70(%rsp), %r8
14008fe82: 4c 89 fa                    	movq	%r15, %rdx
14008fe85: 4c 89 c9                    	movq	%r9, %rcx
14008fe88: e8 a3 e8 ff ff              	callq	0x14008e730 <.text+0x8d730>
14008fe8d: 4c 8b a4 24 80 00 00 00     	movq	0x80(%rsp), %r12
14008fe95: 4c 8b bc 24 88 00 00 00     	movq	0x88(%rsp), %r15
14008fe9d: 31 c0                       	xorl	%eax, %eax
14008fe9f: 48 8d bc 24 e0 00 00 00     	leaq	0xe0(%rsp), %rdi
14008fea7: 4c 8b b4 24 90 00 00 00     	movq	0x90(%rsp), %r14
14008feaf: b9 20 00 00 00              	movl	$0x20, %ecx
14008feb4: f3 aa                       	rep		stosb	%al, %es:(%rdi)
14008feb6: b9 20 00 00 00              	movl	$0x20, %ecx
14008febb: 48 8d bc 24 00 01 00 00     	leaq	0x100(%rsp), %rdi
14008fec3: f3 aa                       	rep		stosb	%al, %es:(%rdi)
14008fec5: b9 20 00 00 00              	movl	$0x20, %ecx
14008feca: 48 8d bc 24 20 01 00 00     	leaq	0x120(%rsp), %rdi
14008fed2: f3 aa                       	rep		stosb	%al, %es:(%rdi)
14008fed4: 48 89 f7                    	movq	%rsi, %rdi
14008fed7: 48 89 d9                    	movq	%rbx, %rcx
14008feda: f3 aa                       	rep		stosb	%al, %es:(%rdi)
14008fedc: e9 fb fe ff ff              	jmp	0x14008fddc <.text+0x8eddc>
14008fee1: 48 b8 2e 31 30 3a 38 37 38 38       	movabsq	$0x383837383a30312e, %rax # imm = 0x383837383A30312E
14008feeb: 48 39 41 10                 	cmpq	%rax, 0x10(%rcx)
14008feef: 0f 85 c7 fe ff ff           	jne	0x14008fdbc <.text+0x8edbc>
14008fef5: 48 8d 84 24 40 01 00 00     	leaq	0x140(%rsp), %rax
14008fefd: 48 8d 15 11 b8 06 00        	leaq	0x6b811(%rip), %rdx     # 0x1400fb715
14008ff04: 48 8d b4 24 50 01 00 00     	leaq	0x150(%rsp), %rsi
14008ff0c: 4c 8d 42 09                 	leaq	0x9(%rdx), %r8
14008ff10: 48 89 c1                    	movq	%rax, %rcx
14008ff13: 48 89 44 24 40              	movq	%rax, 0x40(%rsp)
14008ff18: 48 89 b4 24 40 01 00 00     	movq	%rsi, 0x140(%rsp)
14008ff20: e8 7b e5 ff ff              	callq	0x14008e4a0 <.text+0x8d4a0>
14008ff25: 48 8b 4c 24 40              	movq	0x40(%rsp), %rcx
14008ff2a: e8 a1 ea ff ff              	callq	0x14008e9d0 <.text+0x8d9d0>
14008ff2f: 84 c0                       	testb	%al, %al
14008ff31: 74 2d                       	je	0x14008ff60 <.text+0x8ef60>
14008ff33: 48 8b 8c 24 40 01 00 00     	movq	0x140(%rsp), %rcx
14008ff3b: 31 ed                       	xorl	%ebp, %ebp
14008ff3d: 48 39 f1                    	cmpq	%rsi, %rcx
14008ff40: 0f 84 07 fe ff ff           	je	0x14008fd4d <.text+0x8ed4d>
14008ff46: 48 8b 84 24 50 01 00 00     	movq	0x150(%rsp), %rax
14008ff4e: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008ff52: e8 79 a1 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008ff57: e9 f3 fd ff ff              	jmp	0x14008fd4f <.text+0x8ed4f>
14008ff5c: 0f 1f 40 00                 	nopl	(%rax)
14008ff60: 48 8d 15 b8 b7 06 00        	leaq	0x6b7b8(%rip), %rdx     # 0x1400fb71f
14008ff67: 48 8d bc 24 30 01 00 00     	leaq	0x130(%rsp), %rdi
14008ff6f: 4c 8d 42 0a                 	leaq	0xa(%rdx), %r8
14008ff73: 48 8d 8c 24 20 01 00 00     	leaq	0x120(%rsp), %rcx
14008ff7b: 48 89 bc 24 20 01 00 00     	movq	%rdi, 0x120(%rsp)
14008ff83: e8 18 e5 ff ff              	callq	0x14008e4a0 <.text+0x8d4a0>
14008ff88: 48 8d 8c 24 20 01 00 00     	leaq	0x120(%rsp), %rcx
14008ff90: e8 3b ea ff ff              	callq	0x14008e9d0 <.text+0x8d9d0>
14008ff95: 48 8b 8c 24 20 01 00 00     	movq	0x120(%rsp), %rcx
14008ff9d: 89 c5                       	movl	%eax, %ebp
14008ff9f: 48 39 f9                    	cmpq	%rdi, %rcx
14008ffa2: 74 11                       	je	0x14008ffb5 <.text+0x8efb5>
14008ffa4: 48 8b 84 24 30 01 00 00     	movq	0x130(%rsp), %rax
14008ffac: 48 8d 50 01                 	leaq	0x1(%rax), %rdx
14008ffb0: e8 1b a1 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14008ffb5: 48 8b 8c 24 40 01 00 00     	movq	0x140(%rsp), %rcx
14008ffbd: 83 f5 01                    	xorl	$0x1, %ebp
14008ffc0: 48 39 f1                    	cmpq	%rsi, %rcx
14008ffc3: 75 81                       	jne	0x14008ff46 <.text+0x8ef46>
14008ffc5: e9 85 fd ff ff              	jmp	0x14008fd4f <.text+0x8ed4f>
14008ffca: 48 b8 ff ff ff ff ff ff ff 7f       	movabsq	$0x7fffffffffffffff, %rax # imm = 0x7FFFFFFFFFFFFFFF
14008ffd4: 4c 29 c0                    	subq	%r8, %rax
14008ffd7: 48 39 e8                    	cmpq	%rbp, %rax
14008ffda: 0f 82 87 87 04 00           	jb	0x1400d8767 <.text+0xd7767>
14008ffe0: 4c 39 c5                    	cmpq	%r8, %rbp
14008ffe3: 4c 89 c0                    	movq	%r8, %rax
14008ffe6: 48 0f 43 c5                 	cmovaeq	%rbp, %rax
14008ffea: 4a 8d 0c 00                 	leaq	(%rax,%r8), %rcx
14008ffee: 48 b8 ff ff ff ff ff ff ff 7f       	movabsq	$0x7fffffffffffffff, %rax # imm = 0x7FFFFFFFFFFFFFFF
14008fff8: 48 39 c1                    	cmpq	%rax, %rcx
14008fffb: 48 0f 46 c1                 	cmovbeq	%rcx, %rax
14008ffff: 48 89 c1                    	movq	%rax, %rcx
140090002: 48 89 44 24 50              	movq	%rax, 0x50(%rsp)
140090007: e8 d4 a0 03 00              	callq	0x1400ca0e0 <.text+0xc90e0>
14009000c: 49 89 c1                    	movq	%rax, %r9
14009000f: 4d 39 fc                    	cmpq	%r15, %r12
140090012: 0f 84 7e 00 00 00           	je	0x140090096 <.text+0x8f096>
140090018: 4d 29 e7                    	subq	%r12, %r15
14009001b: 31 c0                       	xorl	%eax, %eax
14009001d: 0f 1f 00                    	nopl	(%rax)
140090020: 41 0f b6 14 04              	movzbl	(%r12,%rax), %edx
140090025: 41 88 14 01                 	movb	%dl, (%r9,%rax)
140090029: 48 83 c0 01                 	addq	$0x1, %rax
14009002d: 49 39 c7                    	cmpq	%rax, %r15
140090030: 75 ee                       	jne	0x140090020 <.text+0x8f020>
140090032: 4d 01 cf                    	addq	%r9, %r15
140090035: 48 8b 54 24 38              	movq	0x38(%rsp), %rdx
14009003a: 49 89 e8                    	movq	%rbp, %r8
14009003d: 4c 89 4c 24 58              	movq	%r9, 0x58(%rsp)
140090042: 4c 89 f9                    	movq	%r15, %rcx
140090045: 49 01 ef                    	addq	%rbp, %r15
140090048: e8 9b 57 01 00              	callq	0x1400a57e8 <.text+0xa47e8>
14009004d: 4c 8b 4c 24 58              	movq	0x58(%rsp), %r9
140090052: 0f 1f 00                    	nopl	(%rax)
140090055: 66 66 2e 0f 1f 84 00 00 00 00 00    	nopw	%cs:(%rax,%rax)
140090060: 4c 89 f2                    	movq	%r14, %rdx
140090063: 4c 89 e1                    	movq	%r12, %rcx
140090066: 4c 89 4c 24 58              	movq	%r9, 0x58(%rsp)
14009006b: 4c 29 e2                    	subq	%r12, %rdx
14009006e: e8 5d a0 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
140090073: 4c 8b 4c 24 58              	movq	0x58(%rsp), %r9
140090078: 0f 1f 84 00 00 00 00 00     	nopl	(%rax,%rax)
140090080: 4c 8b 74 24 50              	movq	0x50(%rsp), %r14
140090085: 4d 89 f8                    	movq	%r15, %r8
140090088: 4d 89 cc                    	movq	%r9, %r12
14009008b: 4d 29 c8                    	subq	%r9, %r8
14009008e: 4d 01 ce                    	addq	%r9, %r14
140090091: e9 c8 fb ff ff              	jmp	0x14008fc5e <.text+0x8ec5e>
140090096: 48 8b 54 24 38              	movq	0x38(%rsp), %rdx
14009009b: 49 89 e8                    	movq	%rbp, %r8
14009009e: 48 89 c1                    	movq	%rax, %rcx
1400900a1: e8 42 57 01 00              	callq	0x1400a57e8 <.text+0xa47e8>
1400900a6: 49 89 c1                    	movq	%rax, %r9
1400900a9: 4c 8d 3c 28                 	leaq	(%rax,%rbp), %r15
1400900ad: 4d 85 e4                    	testq	%r12, %r12
1400900b0: 74 ce                       	je	0x140090080 <.text+0x8f080>
1400900b2: eb ac                       	jmp	0x140090060 <.text+0x8f060>
1400900b4: e9 ca 86 04 00              	jmp	0x1400d8783 <.text+0xd7783>
1400900b9: 48 89 c3                    	movq	%rax, %rbx
1400900bc: e9 08 87 04 00              	jmp	0x1400d87c9 <.text+0xd77c9>
1400900c1: e9 2e 87 04 00              	jmp	0x1400d87f4 <.text+0xd77f4>
1400900c6: e9 3e 87 04 00              	jmp	0x1400d8809 <.text+0xd7809>
1400900cb: 48 89 c3                    	movq	%rax, %rbx
1400900ce: e9 03 87 04 00              	jmp	0x1400d87d6 <.text+0xd77d6>
1400900d3: 48 89 c3                    	movq	%rax, %rbx
1400900d6: e9 0a 87 04 00              	jmp	0x1400d87e5 <.text+0xd77e5>
1400900db: e9 3e 87 04 00              	jmp	0x1400d881e <.text+0xd781e>
1400900e0: e9 4e 87 04 00              	jmp	0x1400d8833 <.text+0xd7833>
1400900e5: e9 53 87 04 00              	jmp	0x1400d883d <.text+0xd783d>
1400900ea: 48 89 c3                    	movq	%rax, %rbx
1400900ed: e9 56 87 04 00              	jmp	0x1400d8848 <.text+0xd7848>
1400900f2: e9 5c 87 04 00              	jmp	0x1400d8853 <.text+0xd7853>
1400900f7: 48 89 c3                    	movq	%rax, %rbx
1400900fa: e9 79 87 04 00              	jmp	0x1400d8878 <.text+0xd7878>
1400900ff: 48 89 c3                    	movq	%rax, %rbx
140090102: e9 de 86 04 00              	jmp	0x1400d87e5 <.text+0xd77e5>
140090107: 48 89 c3                    	movq	%rax, %rbx
14009010a: e9 72 87 04 00              	jmp	0x1400d8881 <.text+0xd7881>
14009010f: 48 89 c3                    	movq	%rax, %rbx
140090112: e9 75 87 04 00              	jmp	0x1400d888c <.text+0xd788c>
140090117: e9 7b 87 04 00              	jmp	0x1400d8897 <.text+0xd7897>
14009011c: e9 52 86 04 00              	jmp	0x1400d8773 <.text+0xd7773>
140090121: e9 5d 86 04 00              	jmp	0x1400d8783 <.text+0xd7783>
140090126: 66 2e 0f 1f 84 00 00 00 00 00       	nopw	%cs:(%rax,%rax)
140090130: 41 57                       	pushq	%r15
140090132: 41 56                       	pushq	%r14
140090134: 41 55                       	pushq	%r13
140090136: 41 54                       	pushq	%r12
140090138: 55                          	pushq	%rbp
140090139: 57                          	pushq	%rdi
14009013a: 56                          	pushq	%rsi
14009013b: 53                          	pushq	%rbx
14009013c: 48 81 ec c8 02 00 00        	subq	$0x2c8, %rsp            # imm = 0x2C8
140090143: 41 b8 08 02 00 00           	movl	$0x208, %r8d            # imm = 0x208
140090149: 31 d2                       	xorl	%edx, %edx
14009014b: 48 89 cb                    	movq	%rcx, %rbx
14009014e: 4c 8d 6c 24 60              	leaq	0x60(%rsp), %r13
140090153: 48 8d 8c 24 b0 00 00 00     	leaq	0xb0(%rsp), %rcx
14009015b: c6 44 24 60 00              	movb	$0x0, 0x60(%rsp)
140090160: 4c 89 6c 24 50              	movq	%r13, 0x50(%rsp)
140090165: 48 8d 74 24 50              	leaq	0x50(%rsp), %rsi
14009016a: 48 c7 44 24 58 00 00 00 00  	movq	$0x0, 0x58(%rsp)
140090173: e8 d8 54 01 00              	callq	0x1400a5650 <.text+0xa4650>
140090178: 41 b8 04 01 00 00           	movl	$0x104, %r8d            # imm = 0x104
14009017e: 48 8d 94 24 b0 00 00 00     	leaq	0xb0(%rsp), %rdx
140090186: 31 c9                       	xorl	%ecx, %ecx
140090188: ff 15 c2 5c 0e 00           	callq	*0xe5cc2(%rip)          # 0x140175e50
14009018e: 48 8d 8c 24 b0 00 00 00     	leaq	0xb0(%rsp), %rcx
140090196: 48 8d bc 24 a0 00 00 00     	leaq	0xa0(%rsp), %rdi
14009019e: 48 89 bc 24 90 00 00 00     	movq	%rdi, 0x90(%rsp)
1400901a6: e8 d5 54 01 00              	callq	0x1400a5680 <.text+0xa4680>
1400901ab: 48 8d 94 24 b0 00 00 00     	leaq	0xb0(%rsp), %rdx
1400901b3: 48 8d 8c 24 90 00 00 00     	leaq	0x90(%rsp), %rcx
1400901bb: 4c 8d 84 44 b0 00 00 00     	leaq	0xb0(%rsp,%rax,2), %r8
1400901c3: e8 98 e3 ff ff              	callq	0x14008e560 <.text+0x8d560>
1400901c8: 48 8d 0d 5b b5 06 00        	leaq	0x6b55b(%rip), %rcx     # 0x1400fb72a
1400901cf: 48 8d b4 24 80 00 00 00     	leaq	0x80(%rsp), %rsi
1400901d7: e8 a4 54 01 00              	callq	0x1400a5680 <.text+0xa4680>
1400901dc: 49 c7 c0 ff ff ff ff        	movq	$-0x1, %r8
1400901e3: 48 8d 6c 24 70              	leaq	0x70(%rsp), %rbp
1400901e8: 48 8d 15 3b b5 06 00        	leaq	0x6b53b(%rip), %rdx     # 0x1400fb72a
1400901ef: 49 89 c1                    	movq	%rax, %r9
1400901f2: 48 8d 8c 24 90 00 00 00     	leaq	0x90(%rsp), %rcx
1400901fa: e8 71 79 02 00              	callq	0x1400b7b70 <.text+0xb6b70>
1400901ff: 48 89 74 24 70              	movq	%rsi, 0x70(%rsp)
140090204: 48 83 f8 ff                 	cmpq	$-0x1, %rax
140090208: 0f 84 32 04 00 00           	je	0x140090640 <.text+0x8f640>
14009020e: 48 8b 8c 24 98 00 00 00     	movq	0x98(%rsp), %rcx
140090216: 48 8b 94 24 90 00 00 00     	movq	0x90(%rsp), %rdx
14009021e: 48 39 c8                    	cmpq	%rcx, %rax
140090221: 48 0f 47 c1                 	cmovaq	%rcx, %rax
140090225: 48 89 e9                    	movq	%rbp, %rcx
140090228: 4c 8d 04 42                 	leaq	(%rdx,%rax,2), %r8
14009022c: e8 2f e3 ff ff              	callq	0x14008e560 <.text+0x8d560>
140090231: 48 8b 8c 24 90 00 00 00     	movq	0x90(%rsp), %rcx
140090239: 48 39 f9                    	cmpq	%rdi, %rcx
14009023c: 74 12                       	je	0x140090250 <.text+0x8f250>
14009023e: 48 8b 84 24 a0 00 00 00     	movq	0xa0(%rsp), %rax
140090246: 48 8d 54 00 02              	leaq	0x2(%rax,%rax), %rdx
14009024b: e8 80 9e 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
140090250: 48 8d 0d dd b4 06 00        	leaq	0x6b4dd(%rip), %rcx     # 0x1400fb734
140090257: e8 24 54 01 00              	callq	0x1400a5680 <.text+0xa4680>
14009025c: 48 ba fe ff ff ff ff ff ff 3f       	movabsq	$0x3ffffffffffffffe, %rdx # imm = 0x3FFFFFFFFFFFFFFE
140090266: 48 2b 54 24 78              	subq	0x78(%rsp), %rdx
14009026b: 49 89 c0                    	movq	%rax, %r8
14009026e: 48 39 c2                    	cmpq	%rax, %rdx
140090271: 0f 82 71 07 00 00           	jb	0x1400909e8 <.text+0x8f9e8>
140090277: 48 8d 15 b6 b4 06 00        	leaq	0x6b4b6(%rip), %rdx     # 0x1400fb734
14009027e: 48 89 e9                    	movq	%rbp, %rcx
140090281: e8 7a 63 03 00              	callq	0x1400c6600 <.text+0xc5600>
140090286: 4c 8d a4 24 c0 00 00 00     	leaq	0xc0(%rsp), %r12
14009028e: 48 8d 50 10                 	leaq	0x10(%rax), %rdx
140090292: 4c 89 a4 24 b0 00 00 00     	movq	%r12, 0xb0(%rsp)
14009029a: 48 8b 08                    	movq	(%rax), %rcx
14009029d: 48 39 d1                    	cmpq	%rdx, %rcx
1400902a0: 0f 84 fa 03 00 00           	je	0x1400906a0 <.text+0x8f6a0>
1400902a6: 48 89 8c 24 b0 00 00 00     	movq	%rcx, 0xb0(%rsp)
1400902ae: 48 8b 48 10                 	movq	0x10(%rax), %rcx
1400902b2: 48 89 8c 24 c0 00 00 00     	movq	%rcx, 0xc0(%rsp)
1400902ba: 48 8b 48 08                 	movq	0x8(%rax), %rcx
1400902be: 45 31 c0                    	xorl	%r8d, %r8d
1400902c1: 48 89 8c 24 b8 00 00 00     	movq	%rcx, 0xb8(%rsp)
1400902c9: 48 8d 0d d0 b4 06 00        	leaq	0x6b4d0(%rip), %rcx     # 0x1400fb7a0
1400902d0: 48 89 10                    	movq	%rdx, (%rax)
1400902d3: 66 44 89 40 10              	movw	%r8w, 0x10(%rax)
1400902d8: 48 c7 40 08 00 00 00 00     	movq	$0x0, 0x8(%rax)
1400902e0: e8 9b 53 01 00              	callq	0x1400a5680 <.text+0xa4680>
1400902e5: 48 ba fe ff ff ff ff ff ff 3f       	movabsq	$0x3ffffffffffffffe, %rdx # imm = 0x3FFFFFFFFFFFFFFE
1400902ef: 48 2b 94 24 b8 00 00 00     	subq	0xb8(%rsp), %rdx
1400902f7: 49 89 c0                    	movq	%rax, %r8
1400902fa: 48 39 c2                    	cmpq	%rax, %rdx
1400902fd: 0f 82 ea 85 04 00           	jb	0x1400d88ed <.text+0xd78ed>
140090303: 48 8d 15 96 b4 06 00        	leaq	0x6b496(%rip), %rdx     # 0x1400fb7a0
14009030a: 48 8d 8c 24 b0 00 00 00     	leaq	0xb0(%rsp), %rcx
140090312: e8 e9 62 03 00              	callq	0x1400c6600 <.text+0xc5600>
140090317: 48 89 bc 24 90 00 00 00     	movq	%rdi, 0x90(%rsp)
14009031f: 48 8b 08                    	movq	(%rax), %rcx
140090322: 48 8d 50 10                 	leaq	0x10(%rax), %rdx
140090326: 48 39 d1                    	cmpq	%rdx, %rcx
140090329: 0f 84 39 03 00 00           	je	0x140090668 <.text+0x8f668>
14009032f: 48 89 8c 24 90 00 00 00     	movq	%rcx, 0x90(%rsp)
140090337: 48 8b 48 10                 	movq	0x10(%rax), %rcx
14009033b: 48 89 8c 24 a0 00 00 00     	movq	%rcx, 0xa0(%rsp)
140090343: 48 8b 48 08                 	movq	0x8(%rax), %rcx
140090347: 48 89 8c 24 98 00 00 00     	movq	%rcx, 0x98(%rsp)
14009034f: 31 c9                       	xorl	%ecx, %ecx
140090351: 48 89 10                    	movq	%rdx, (%rax)
140090354: 66 89 48 10                 	movw	%cx, 0x10(%rax)
140090358: 48 8b 8c 24 b0 00 00 00     	movq	0xb0(%rsp), %rcx
140090360: 48 c7 40 08 00 00 00 00     	movq	$0x0, 0x8(%rax)
140090368: 4c 39 e1                    	cmpq	%r12, %rcx
14009036b: 74 12                       	je	0x14009037f <.text+0x8f37f>
14009036d: 48 8b 84 24 c0 00 00 00     	movq	0xc0(%rsp), %rax
140090375: 48 8d 54 00 02              	leaq	0x2(%rax,%rax), %rdx
14009037a: e8 51 9d 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14009037f: 48 8b 4c 24 70              	movq	0x70(%rsp), %rcx
140090384: 48 39 f1                    	cmpq	%rsi, %rcx
140090387: 74 12                       	je	0x14009039b <.text+0x8f39b>
140090389: 48 8b 84 24 80 00 00 00     	movq	0x80(%rsp), %rax
140090391: 48 8d 54 00 02              	leaq	0x2(%rax,%rax), %rdx
140090396: e8 35 9d 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
14009039b: c7 44 24 28 00 00 00 00     	movl	$0x0, 0x28(%rsp)
1400903a3: 45 31 c9                    	xorl	%r9d, %r9d
1400903a6: 48 8b 8c 24 90 00 00 00     	movq	0x90(%rsp), %rcx
1400903ae: 41 b8 01 00 00 00           	movl	$0x1, %r8d
1400903b4: 48 c7 44 24 30 00 00 00 00  	movq	$0x0, 0x30(%rsp)
1400903bd: ba 00 00 00 80              	movl	$0x80000000, %edx       # imm = 0x80000000
1400903c2: c7 44 24 20 03 00 00 00     	movl	$0x3, 0x20(%rsp)
1400903ca: ff 15 78 59 0e 00           	callq	*0xe5978(%rip)          # 0x140175d48
1400903d0: 48 8b 8c 24 90 00 00 00     	movq	0x90(%rsp), %rcx
1400903d8: 48 89 c5                    	movq	%rax, %rbp
1400903db: 48 39 f9                    	cmpq	%rdi, %rcx
1400903de: 74 12                       	je	0x1400903f2 <.text+0x8f3f2>
1400903e0: 48 8b 84 24 a0 00 00 00     	movq	0xa0(%rsp), %rax
1400903e8: 48 8d 54 00 02              	leaq	0x2(%rax,%rax), %rdx
1400903ed: e8 de 9c 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
1400903f2: 48 83 fd ff                 	cmpq	$-0x1, %rbp
1400903f6: 0f 84 c5 00 00 00           	je	0x1400904c1 <.text+0x8f4c1>
1400903fc: 31 d2                       	xorl	%edx, %edx
1400903fe: 48 89 e9                    	movq	%rbp, %rcx
140090401: 48 8d 74 24 50              	leaq	0x50(%rsp), %rsi
140090406: ff 15 2c 5a 0e 00           	callq	*0xe5a2c(%rip)          # 0x140175e38
14009040c: 41 89 c7                    	movl	%eax, %r15d
14009040f: 8d 40 c0                    	leal	-0x40(%rax), %eax
140090412: 3d c0 0f 00 00              	cmpl	$0xfc0, %eax            # imm = 0xFC0
140090417: 0f 87 13 02 00 00           	ja	0x140090630 <.text+0x8f630>
14009041d: 45 89 fe                    	movl	%r15d, %r14d
140090420: 48 8d 74 24 50              	leaq	0x50(%rsp), %rsi
140090425: 4c 89 f1                    	movq	%r14, %rcx
140090428: e8 b3 9c 03 00              	callq	0x1400ca0e0 <.text+0xc90e0>
14009042d: c6 00 00                    	movb	$0x0, (%rax)
140090430: 4a 8d 34 30                 	leaq	(%rax,%r14), %rsi
140090434: 48 89 c7                    	movq	%rax, %rdi
140090437: 48 8d 48 01                 	leaq	0x1(%rax), %rcx
14009043b: 4d 8d 46 ff                 	leaq	-0x1(%r14), %r8
14009043f: 31 d2                       	xorl	%edx, %edx
140090441: 48 89 74 24 48              	movq	%rsi, 0x48(%rsp)
140090446: 48 89 b4 24 c0 00 00 00     	movq	%rsi, 0xc0(%rsp)
14009044e: 48 89 84 24 b0 00 00 00     	movq	%rax, 0xb0(%rsp)
140090456: e8 f5 51 01 00              	callq	0x1400a5650 <.text+0xa4650>
14009045b: 45 89 f8                    	movl	%r15d, %r8d
14009045e: 48 89 fa                    	movq	%rdi, %rdx
140090461: 48 89 e9                    	movq	%rbp, %rcx
140090464: 48 89 b4 24 b8 00 00 00     	movq	%rsi, 0xb8(%rsp)
14009046c: 4c 8d 8c 24 90 00 00 00     	leaq	0x90(%rsp), %r9
140090474: 48 8d 74 24 50              	leaq	0x50(%rsp), %rsi
140090479: c7 84 24 90 00 00 00 00 00 00 00    	movl	$0x0, 0x90(%rsp)
140090484: 48 c7 44 24 20 00 00 00 00  	movq	$0x0, 0x20(%rsp)
14009048d: ff 15 ad 5a 0e 00           	callq	*0xe5aad(%rip)          # 0x140175f40
140090493: 89 44 24 44                 	movl	%eax, 0x44(%rsp)
140090497: 48 89 e9                    	movq	%rbp, %rcx
14009049a: ff 15 80 58 0e 00           	callq	*0xe5880(%rip)          # 0x140175d20
1400904a0: 8b 54 24 44                 	movl	0x44(%rsp), %edx
1400904a4: 85 d2                       	testl	%edx, %edx
1400904a6: 74 0e                       	je	0x1400904b6 <.text+0x8f4b6>
1400904a8: 44 3b bc 24 90 00 00 00     	cmpl	0x90(%rsp), %r15d
1400904b0: 0f 84 ca 02 00 00           	je	0x140090780 <.text+0x8f780>
1400904b6: 4c 89 f2                    	movq	%r14, %rdx
1400904b9: 48 89 f9                    	movq	%rdi, %rcx
1400904bc: e8 0f 9c 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
1400904c1: b9 15 00 00 00              	movl	$0x15, %ecx
1400904c6: 48 8d 74 24 50              	leaq	0x50(%rsp), %rsi
1400904cb: e8 10 9c 03 00              	callq	0x1400ca0e0 <.text+0xc90e0>
1400904d0: c6 00 68                    	movb	$0x68, (%rax)
1400904d3: 66 0f ef c0                 	pxor	%xmm0, %xmm0
1400904d7: 48 89 c6                    	movq	%rax, %rsi
1400904da: b9 01 00 00 00              	movl	$0x1, %ecx
1400904df: c7 40 11 00 00 00 00        	movl	$0x0, 0x11(%rax)
1400904e6: 4c 8d 15 78 b2 06 00        	leaq	0x6b278(%rip), %r10     # 0x1400fb765
1400904ed: 4c 8d 0d 5c b2 06 00        	leaq	0x6b25c(%rip), %r9      # 0x1400fb750
1400904f4: 49 b8 25 49 92 24 49 92 24 49       	movabsq	$0x4924924924924925, %r8 # imm = 0x4924924924924925
1400904fe: 0f 11 40 01                 	movups	%xmm0, 0x1(%rax)
140090502: 0f 1f 80 00 00 00 00        	nopl	(%rax)
140090509: 66 66 2e 0f 1f 84 00 00 00 00 00    	nopw	%cs:(%rax,%rax)
140090514: 66 66 2e 0f 1f 84 00 00 00 00 00    	nopw	%cs:(%rax,%rax)
14009051f: 66 66 2e 0f 1f 84 00 00 00 00 00    	nopw	%cs:(%rax,%rax)
14009052a: 66 66 2e 0f 1f 84 00 00 00 00 00    	nopw	%cs:(%rax,%rax)
140090535: 66 66 2e 0f 1f 84 00 00 00 00 00    	nopw	%cs:(%rax,%rax)
140090540: 48 89 c8                    	movq	%rcx, %rax
140090543: 49 f7 e8                    	imulq	%r8
140090546: 48 89 c8                    	movq	%rcx, %rax
140090549: 48 c1 f8 3f                 	sarq	$0x3f, %rax
14009054d: 48 d1 fa                    	sarq	%rdx
140090550: 48 29 c2                    	subq	%rax, %rdx
140090553: 48 8d 04 d5 00 00 00 00     	leaq	(,%rdx,8), %rax
14009055b: 48 29 d0                    	subq	%rdx, %rax
14009055e: 48 89 ca                    	movq	%rcx, %rdx
140090561: 48 29 c2                    	subq	%rax, %rdx
140090564: 41 0f b6 04 12              	movzbl	(%r10,%rdx), %eax
140090569: 41 32 04 09                 	xorb	(%r9,%rcx), %al
14009056d: 88 04 0e                    	movb	%al, (%rsi,%rcx)
140090570: 48 83 c1 01                 	addq	$0x1, %rcx
140090574: 48 83 f9 15                 	cmpq	$0x15, %rcx
140090578: 75 c6                       	jne	0x140090540 <.text+0x8f540>
14009057a: 48 89 f2                    	movq	%rsi, %rdx
14009057d: 48 8d 7e 15                 	leaq	0x15(%rsi), %rdi
140090581: 48 b9 25 23 22 84 e4 9c f2 cb       	movabsq	$-0x340d631b7bdddcdb, %rcx # imm = 0xCBF29CE484222325
14009058b: 49 b8 b3 01 00 00 00 01 00 00       	movabsq	$0x100000001b3, %r8 # imm = 0x100000001B3
140090595: 66 66 2e 0f 1f 84 00 00 00 00 00    	nopw	%cs:(%rax,%rax)
1400905a0: 0f b6 02                    	movzbl	(%rdx), %eax
1400905a3: 48 83 c2 01                 	addq	$0x1, %rdx
1400905a7: 48 31 c8                    	xorq	%rcx, %rax
1400905aa: 48 89 c1                    	movq	%rax, %rcx
1400905ad: 49 0f af c8                 	imulq	%r8, %rcx
1400905b1: 48 39 fa                    	cmpq	%rdi, %rdx
1400905b4: 75 ea                       	jne	0x1400905a0 <.text+0x8f5a0>
1400905b6: 48 ba c6 65 9c 3a c0 cc 6a 28       	movabsq	$0x286accc03a9c65c6, %rdx # imm = 0x286ACCC03A9C65C6
1400905c0: 48 39 d0                    	cmpq	%rdx, %rax
1400905c3: 0f 84 07 01 00 00           	je	0x1400906d0 <.text+0x8f6d0>
1400905c9: 48 89 f7                    	movq	%rsi, %rdi
1400905cc: b9 15 00 00 00              	movl	$0x15, %ecx
1400905d1: 31 c0                       	xorl	%eax, %eax
1400905d3: f3 aa                       	rep		stosb	%al, %es:(%rdi)
1400905d5: 48 8d 43 10                 	leaq	0x10(%rbx), %rax
1400905d9: 48 c7 43 08 00 00 00 00     	movq	$0x0, 0x8(%rbx)
1400905e1: 48 89 03                    	movq	%rax, (%rbx)
1400905e4: 31 c0                       	xorl	%eax, %eax
1400905e6: 66 89 43 10                 	movw	%ax, 0x10(%rbx)
1400905ea: ba 15 00 00 00              	movl	$0x15, %edx
1400905ef: 48 89 f1                    	movq	%rsi, %rcx
1400905f2: e8 d9 9a 03 00              	callq	0x1400ca0d0 <.text+0xc90d0>
1400905f7: 48 8b 4c 24 50              	movq	0x50(%rsp), %rcx
1400905fc: 4c 39 e9                    	cmpq	%r13, %rcx
1400905ff: 74 0e                       	je	0x14009060f <.text+0x8f60f>
