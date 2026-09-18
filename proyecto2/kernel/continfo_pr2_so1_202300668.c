/*
 * Ejemplo base de módulo de kernel en C para Proyecto 2 SO1.
 * Este archivo sirve como plantilla para extraer métricas de memoria y procesos
 * relacionados con contenedores y escribirlas en /proc/continfo_pr2_so1_202300668.
 *
 * Requiere headers del kernel y compilación con una versión compatible de Linux.
 */

#include <linux/init.h>
#include <linux/kernel.h>
#include <linux/module.h>
#include <linux/proc_fs.h>
#include <linux/seq_file.h>
#include <linux/sched.h>
#include <linux/mm.h>
#include <linux/utsname.h>

MODULE_LICENSE("GPL");
MODULE_AUTHOR("202300668");
MODULE_DESCRIPTION("Kernel telemetry for SO1 project 2");

static const char *proc_name = "continfo_pr2_so1_202300668";
static struct proc_dir_entry *proc_entry;

static int proc_show(struct seq_file *m, void *v)
{
	struct sysinfo info;
	si_meminfo(&info);

	seq_printf(m, "# PID,NAME,COMMAND,VSZ,RSS,MEM_PERCENT,CPU_PERCENT,CONTAINER_ID\n");
	seq_printf(m, "0,ram_total,%lu,0,0,0,0,system\n", (unsigned long)info.totalram * info.mem_unit / 1024);
	seq_printf(m, "0,ram_free,%lu,0,0,0,0,system\n", (unsigned long)info.freeram * info.mem_unit / 1024);
	seq_printf(m, "0,ram_used,%lu,0,0,0,0,system\n", ((unsigned long)info.totalram - (unsigned long)info.freeram) * info.mem_unit / 1024);
	return 0;
}

static int proc_open(struct inode *inode, struct file *file)
{
	return single_open(file, proc_show, NULL);
}

static const struct proc_ops proc_fops = {
	.proc_open = proc_open,
	.proc_read = seq_read,
	.proc_lseek = seq_lseek,
	.proc_release = single_release,
};

static int __init continfo_init(void)
{
	proc_entry = proc_create(proc_name, 0, NULL, &proc_fops);
	if (!proc_entry) {
		pr_err("No se pudo crear /proc/%s\n", proc_name);
		return -ENOMEM;
	}
	pr_info("Modulo cargado: %s\n", proc_name);
	return 0;
}

static void __exit continfo_exit(void)
{
	remove_proc_entry(proc_name, NULL);
	pr_info("Modulo descargado: %s\n", proc_name);
}

module_init(continfo_init);
module_exit(continfo_exit);
